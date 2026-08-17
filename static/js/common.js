window.isPreviewActive = false;
let currentIconUrls = [];
let currentIconIndex = 0;

function isVideoUrl(url) {
    return /\.(mp4|webm)(\?|#|$)/i.test(url);
}

function createMediaElement(src, isMain) {
    const isVideo = isVideoUrl(src);
    if (isVideo) {
        const video = document.createElement('video');
        video.src = src;
        if (isMain) {
            video.style.cssText = `
                max-width: 85vw;
                max-height: 85vh;
                min-width: 100px;
                min-height: 100px;
                flex-shrink: 1;
                border-radius: 4px;
            `;
            video.controls = true;
            video.autoplay = true;
            video.loop = true;
        } else {
            video.style.cssText = `
                width: 15vw;
                height: 15vw;
                max-height: 50vh;
                min-width: 60px;
                min-height: 60px;
                object-fit: cover;
                opacity: 0.6;
                cursor: pointer;
                transition: opacity 0.3s, transform 0.3s;
                flex-shrink: 0;
                border-radius: 4px;
            `;
            video.muted = true;
            video.autoplay = true;
            video.loop = true;
            video.playsInline = true;
            video.onmouseenter = function() {
                this.style.opacity = '1';
                this.style.transform = 'scale(1.1)';
            };
            video.onmouseleave = function() {
                this.style.opacity = '0.6';
                this.style.transform = 'scale(1)';
            };
        }
        return video;
    } else {
        const img = document.createElement('img');
        img.src = src;
        if (isMain) {
            img.style.cssText = `
                max-width: 85vw;
                max-height: 85vh;
                min-width: 100px;
                min-height: 100px;
                object-fit: contain;
                flex-shrink: 1;
            `;
        } else {
            img.style.cssText = `
                width: 15vw;
                height: 15vw;
                max-height: 50vh;
                min-width: 60px;
                min-height: 60px;
                object-fit: contain;
                opacity: 0.6;
                cursor: pointer;
                transition: opacity 0.3s, transform 0.3s;
                flex-shrink: 0;
            `;
            img.onmouseenter = function() {
                this.style.opacity = '1';
                this.style.transform = 'scale(1.1)';
            };
            img.onmouseleave = function() {
                this.style.opacity = '0.6';
                this.style.transform = 'scale(1)';
            };
        }
        return img;
    }
}

function showIconPreview(iconUrls, index = 0) {
    window.isPreviewActive = true;
    currentIconUrls = iconUrls;
    currentIconIndex = index;

    let lightbox = document.getElementById('lightbox');
    if (!lightbox) {
        lightbox = document.createElement('div');
        lightbox.id = 'lightbox';
        lightbox.style.cssText = `
            display: none;
            position: fixed;
            inset: 0;
            width: 100vw;
            height: 100vh;
            background: rgba(0, 0, 0, 0.85);
            z-index: 9999999;
            align-items: center;
            justify-content: center;
            user-select: none;
            touch-action: none;
        `;

        const container = document.createElement('div');
        container.id = 'lightboxContainer';
        container.style.cssText = `
            display: flex;
            align-items: center;
            gap: 20px;
            padding: 20px;
            max-width: 100vw;
            max-height: 90vh;
            box-sizing: border-box;
        `;

        const prevImg = createMediaElement('', false);
        prevImg.id = 'lightboxPrev';

        const mainContainer = document.createElement('div');
        mainContainer.id = 'lightboxMain';
        mainContainer.style.cssText = `
            display: flex;
            align-items: center;
            justify-content: center;
            flex-shrink: 1;
            max-width: 85vw;
            max-height: 85vh;
        `;

        const nextImg = createMediaElement('', false);
        nextImg.id = 'lightboxNext';

        container.appendChild(prevImg);
        container.appendChild(mainContainer);
        container.appendChild(nextImg);
        lightbox.appendChild(container);
        document.body.appendChild(lightbox);

        setupTouchEvents(lightbox);
    }

    updatePreviewImages(currentIconIndex);

    lightbox.onclick = function() {
        closePreview();
    };

    lightbox.style.display = 'flex';
}

function closePreview() {
    const lightbox = document.getElementById('lightbox');
    if (lightbox) {
        // 停止所有视频播放
        const videos = lightbox.querySelectorAll('video');
        videos.forEach(v => { v.pause(); v.src = ''; });
        lightbox.style.display = 'none';
    }
    window.isPreviewActive = false;
}

function setupTouchEvents(lightbox) {
    let touchStartX = 0;
    let touchEndX = 0;
    const minSwipeDistance = 50;

    lightbox.addEventListener('touchstart', function(e) {
        touchStartX = e.touches[0].clientX;
    }, { passive: true });

    lightbox.addEventListener('touchend', function(e) {
        touchEndX = e.changedTouches[0].clientX;
        handleSwipe();
    }, { passive: true });

    function handleSwipe() {
        const distance = touchStartX - touchEndX;
        const isLeftSwipe = distance > minSwipeDistance;
        const isRightSwipe = distance < -minSwipeDistance;

        if (isLeftSwipe) {
            nextIcon();
        } else if (isRightSwipe) {
            prevIcon();
        }
    }
}

function prevIcon() {
    if (!window.isPreviewActive || currentIconUrls.length === 0) return;
    const newIndex = currentIconIndex > 0 ? currentIconIndex - 1 : currentIconUrls.length - 1;
    updatePreviewImages(newIndex);
}

function nextIcon() {
    if (!window.isPreviewActive || currentIconUrls.length === 0) return;
    const newIndex = currentIconIndex < currentIconUrls.length - 1 ? currentIconIndex + 1 : 0;
    updatePreviewImages(newIndex);
}

function updatePreviewImages(index) {
    const mainContainer = document.getElementById('lightboxMain');
    const prevImg = document.getElementById('lightboxPrev');
    const nextImg = document.getElementById('lightboxNext');

    if (!mainContainer || !prevImg || !nextImg) return;

    currentIconIndex = index;

    // 更新主预览
    mainContainer.innerHTML = '';
    const mainEl = createMediaElement(currentIconUrls[currentIconIndex], true);
    mainEl.onclick = function(e) {
        e.stopPropagation();
        closePreview();
    };
    mainContainer.appendChild(mainEl);

    // 更新缩略图
    const prevIndex = currentIconIndex > 0 ? currentIconIndex - 1 : currentIconUrls.length - 1;
    const nextIndex = currentIconIndex < currentIconUrls.length - 1 ? currentIconIndex + 1 : 0;

    const newPrev = createMediaElement(currentIconUrls[prevIndex], false);
    newPrev.id = 'lightboxPrev';
    newPrev.onclick = function(e) {
        e.stopPropagation();
        const newIdx = currentIconIndex > 0 ? currentIconIndex - 1 : currentIconUrls.length - 1;
        updatePreviewImages(newIdx);
    };
    prevImg.replaceWith(newPrev);

    const newNext = createMediaElement(currentIconUrls[nextIndex], false);
    newNext.id = 'lightboxNext';
    newNext.onclick = function(e) {
        e.stopPropagation();
        const newIdx = currentIconIndex < currentIconUrls.length - 1 ? currentIconIndex + 1 : 0;
        updatePreviewImages(newIdx);
    };
    nextImg.replaceWith(newNext);
}

document.addEventListener('keydown', function(e) {
    if (!window.isPreviewActive) return;

    if (e.key === 'ArrowLeft' || e.key === 'ArrowRight') {
        e.preventDefault();
        if (e.key === 'ArrowLeft') {
            prevIcon();
        } else if (e.key === 'ArrowRight') {
            nextIcon();
        }
    }

    if (e.key === 'Escape') {
        closePreview();
    }
});
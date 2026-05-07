window.isPreviewActive = false;
let currentIconUrls = [];
let currentIconIndex = 0;

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
        container.style.cssText = `
            display: flex;
            align-items: center;
            gap: 20px;
            padding: 15px;
            max-width: 95vw;
            max-height: 90vh;
            box-sizing: border-box;
        `;

        const prevImg = document.createElement('img');
        prevImg.id = 'lightboxPrev';
        prevImg.style.cssText = `
            width: 120px;
            height: 120px;
            max-width: 20vw;
            max-height: 20vh;
            min-width: 50px;
            min-height: 50px;
            object-fit: contain;
            opacity: 0.6;
            cursor: pointer;
            transition: opacity 0.3s, transform 0.3s;
            flex-shrink: 0;
        `;
        prevImg.onmouseenter = function() {
            this.style.opacity = '1';
            this.style.transform = 'scale(1.1)';
        };
        prevImg.onmouseleave = function() {
            this.style.opacity = '0.6';
            this.style.transform = 'scale(1)';
        };

        const mainImg = document.createElement('img');
        mainImg.id = 'lightboxMain';
        mainImg.style.cssText = `
            width: 300px;
            height: 300px;
            max-width: 60vw;
            max-height: 70vh;
            min-width: 80px;
            min-height: 80px;
            object-fit: contain;
            flex-shrink: 1;
        `;

        const nextImg = document.createElement('img');
        nextImg.id = 'lightboxNext';
        nextImg.style.cssText = `
            width: 120px;
            height: 120px;
            max-width: 20vw;
            max-height: 20vh;
            min-width: 50px;
            min-height: 50px;
            object-fit: contain;
            opacity: 0.6;
            cursor: pointer;
            transition: opacity 0.3s, transform 0.3s;
            flex-shrink: 0;
        `;
        nextImg.onmouseenter = function() {
            this.style.opacity = '1';
            this.style.transform = 'scale(1.1)';
        };
        nextImg.onmouseleave = function() {
            this.style.opacity = '0.6';
            this.style.transform = 'scale(1)';
        };

        container.appendChild(prevImg);
        container.appendChild(mainImg);
        container.appendChild(nextImg);
        lightbox.appendChild(container);
        document.body.appendChild(lightbox);

        setupTouchEvents(lightbox);
    }

    const prevImg = document.getElementById('lightboxPrev');
    const mainImg = document.getElementById('lightboxMain');
    const nextImg = document.getElementById('lightboxNext');

    const updateImages = (index) => {
        currentIconIndex = index;
        mainImg.src = currentIconUrls[currentIconIndex];

        const prevIndex = currentIconIndex > 0 ? currentIconIndex - 1 : currentIconUrls.length - 1;
        const nextIndex = currentIconIndex < currentIconUrls.length - 1 ? currentIconIndex + 1 : 0;

        prevImg.src = currentIconUrls[prevIndex];
        nextImg.src = currentIconUrls[nextIndex];
    };

    updateImages(currentIconIndex);

    prevImg.onclick = function(e) {
        e.stopPropagation();
        const newIndex = currentIconIndex > 0 ? currentIconIndex - 1 : currentIconUrls.length - 1;
        updateImages(newIndex);
    };

    nextImg.onclick = function(e) {
        e.stopPropagation();
        const newIndex = currentIconIndex < currentIconUrls.length - 1 ? currentIconIndex + 1 : 0;
        updateImages(newIndex);
    };

    lightbox.onclick = function() {
        closePreview();
    };

    mainImg.onclick = function(e) {
        e.stopPropagation();
    };

    lightbox.style.display = 'flex';
}

function closePreview() {
    const lightbox = document.getElementById('lightbox');
    if (lightbox) {
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
    const prevImg = document.getElementById('lightboxPrev');
    const mainImg = document.getElementById('lightboxMain');
    const nextImg = document.getElementById('lightboxNext');

    if (!mainImg || !prevImg || !nextImg) return;

    currentIconIndex = index;
    mainImg.src = currentIconUrls[currentIconIndex];

    const prevIndex = currentIconIndex > 0 ? currentIconIndex - 1 : currentIconUrls.length - 1;
    const nextIndex = currentIconIndex < currentIconUrls.length - 1 ? currentIconIndex + 1 : 0;

    prevImg.src = currentIconUrls[prevIndex];
    nextImg.src = currentIconUrls[nextIndex];
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
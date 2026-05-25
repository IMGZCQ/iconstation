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
            padding: 20px;
            max-width: 100vw;
            max-height: 90vh;
            box-sizing: border-box;
        `;

        const prevImg = document.createElement('img');
        prevImg.id = 'lightboxPrev';
        prevImg.style.cssText = `
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
            width: 60vw;
            height: 60vw;
            max-height: 85vh;
            min-width: 100px;
            min-height: 100px;
            object-fit: contain;
            flex-shrink: 1;
        `;

        const nextImg = document.createElement('img');
        nextImg.id = 'lightboxNext';
        nextImg.style.cssText = `
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
        closePreview();
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

function addParentWindowButtons() {
    try {
        if (window.self === window.top) {
            console.log('不在 iframe 中运行，跳过添加按钮');
            return;
        }

        const parentDoc = window.parent.document;

        const currentIframe = parentDoc.querySelector(`iframe[src="${window.location.href}"]`) ||
                             [...parentDoc.querySelectorAll('iframe')].find(iframe => {
                                 try {
                                     return iframe.contentWindow === window;
                                 } catch (e) {
                                     return false;
                                 }
                             });

        if (!currentIframe) {
            console.error('未找到当前 iframe 元素');
            return;
        }

        const headerContainer = currentIframe.closest('.trim-ui__app-layout--header');
        if (!headerContainer) {
            console.error('未找到包含当前 iframe 的 header 容器');
            return;
        }

        const buttonContainer = headerContainer.querySelector(':scope > div:last-child');
        if (!buttonContainer || !buttonContainer.classList.contains('items-center')) {
            console.error('未找到按钮容器');
            return;
        }

        if (parentDoc.getElementById('iframe-refresh-btn')) {
            console.log('按钮已存在，跳过重复添加');
            return;
        }

        const refreshBtn = createButton(parentDoc, 'iframe-refresh-btn', '刷新页面',
            `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"/>
            </svg>`,
            () => { window.location.reload(); },
            (el) => { el.style.transform = 'rotate(90deg)'; },
            (el) => { el.style.transform = 'rotate(0deg)'; }
        );

        const openNewWindowBtn = createButton(parentDoc, 'iframe-open-new-window-btn', '新标签页打开',
            `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24">
                <path d="M0 0h24v24H0z" fill="none" />
                <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 4H6a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-4m-8-2l8-8m0 0v5m0-5h-5" />
            </svg>`,
            () => { window.open(window.location.href, '_blank', 'noopener'); },
            (el) => { el.style.transform = 'scale(1.1)'; },
            (el) => { el.style.transform = 'scale(1)'; }
        );

        buttonContainer.insertBefore(refreshBtn, buttonContainer.firstChild);
        buttonContainer.insertBefore(openNewWindowBtn, refreshBtn.nextSibling);

        console.log('✅ 已成功在父级窗口标题栏添加按钮（仅限当前 iframe 所在的窗口）');
    } catch (error) {
        if (error.message.includes('cross-origin') || error.message.includes('permission')) {
            console.log('跨域限制：无法访问父级窗口（这是正常的安全行为）');
        } else {
            console.error('添加父级窗口按钮失败:', error);
        }
    }
}

function createButton(parentDoc, id, title, svgHTML, onClick, onEnter, onLeave) {
    const btn = parentDoc.createElement('div');
    btn.id = id;
    btn.title = title;
    btn.className = 'flex h-full w-base shrink-0 cursor-pointer items-center justify-center px-[15px] text-[var(--semi-color-text-0)] hover:bg-[var(--semi-color-fill-0)] active:bg-[var(--semi-color-fill-0)]';
    btn.style.transition = 'all 0.2s ease';
    btn.innerHTML = svgHTML;

    btn.onmouseenter = function() {
        this.style.backgroundColor = 'var(--semi-color-fill-0)';
        onEnter(this);
    };
    btn.onmouseleave = function() {
        this.style.backgroundColor = '';
        onLeave(this);
    };
    btn.onclick = function(e) {
        e.stopPropagation();
        onClick();
    };

    return btn;
}

if (window.self !== window.top) {
    addParentWindowButtons();
}

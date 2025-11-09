/**
 * Pedantic Raven - GitHub Pages JavaScript
 * Minimal interactivity and enhancements
 */

(function() {
    'use strict';

    /**
     * Smooth scrolling for anchor links
     */
    function initSmoothScroll() {
        document.querySelectorAll('a[href^="#"]').forEach(anchor => {
            anchor.addEventListener('click', function(e) {
                const href = this.getAttribute('href');
                if (href === '#') return;

                e.preventDefault();
                const target = document.querySelector(href);
                if (target) {
                    target.scrollIntoView({
                        behavior: 'smooth',
                        block: 'start'
                    });

                    // Update URL without triggering scroll
                    history.pushState(null, null, href);
                }
            });
        });
    }

    /**
     * Add active class to current nav link
     */
    function initActiveNav() {
        const sections = document.querySelectorAll('.section[id]');
        const navLinks = document.querySelectorAll('.nav-links a[href^="#"]');

        function updateActiveNav() {
            let current = '';

            sections.forEach(section => {
                const sectionTop = section.offsetTop;
                const sectionHeight = section.clientHeight;
                if (window.pageYOffset >= sectionTop - 100) {
                    current = section.getAttribute('id');
                }
            });

            navLinks.forEach(link => {
                link.classList.remove('active');
                if (link.getAttribute('href') === `#${current}`) {
                    link.classList.add('active');
                }
            });
        }

        window.addEventListener('scroll', updateActiveNav);
        updateActiveNav();
    }

    /**
     * Lazy load images (diagrams)
     */
    function initLazyLoad() {
        if ('IntersectionObserver' in window) {
            const imageObserver = new IntersectionObserver((entries, observer) => {
                entries.forEach(entry => {
                    if (entry.isIntersecting) {
                        const img = entry.target;
                        if (img.dataset.src) {
                            img.src = img.dataset.src;
                            img.removeAttribute('data-src');
                        }
                        observer.unobserve(img);
                    }
                });
            });

            document.querySelectorAll('img[data-src]').forEach(img => {
                imageObserver.observe(img);
            });
        } else {
            // Fallback for browsers without IntersectionObserver
            document.querySelectorAll('img[data-src]').forEach(img => {
                img.src = img.dataset.src;
                img.removeAttribute('data-src');
            });
        }
    }

    /**
     * Add copy button to code blocks
     */
    function initCodeCopy() {
        document.querySelectorAll('pre code').forEach(block => {
            const pre = block.parentElement;
            const button = document.createElement('button');
            button.className = 'copy-btn';
            button.textContent = 'Copy';
            button.setAttribute('aria-label', 'Copy code to clipboard');

            button.addEventListener('click', async () => {
                try {
                    await navigator.clipboard.writeText(block.textContent);
                    button.textContent = 'Copied!';
                    button.classList.add('copied');

                    setTimeout(() => {
                        button.textContent = 'Copy';
                        button.classList.remove('copied');
                    }, 2000);
                } catch (err) {
                    console.error('Failed to copy:', err);
                    button.textContent = 'Failed';
                    setTimeout(() => {
                        button.textContent = 'Copy';
                    }, 2000);
                }
            });

            // Wrap pre in container for positioning
            const container = document.createElement('div');
            container.className = 'code-container';
            pre.parentNode.insertBefore(container, pre);
            container.appendChild(pre);
            container.appendChild(button);
        });
    }

    /**
     * Add styles for copy button
     */
    function addCopyButtonStyles() {
        const style = document.createElement('style');
        style.textContent = `
            .code-container {
                position: relative;
                margin-bottom: var(--spacing-md, 1rem);
            }

            .copy-btn {
                position: absolute;
                top: 0.5rem;
                right: 0.5rem;
                background-color: var(--raven-slate, #2C3354);
                color: var(--raven-white, #F5F6FA);
                border: 1px solid var(--raven-gray, #4A5172);
                padding: 0.4em 0.8em;
                font-size: 0.85rem;
                border-radius: 4px;
                cursor: pointer;
                transition: all 150ms ease;
                font-family: var(--font-mono, monospace);
            }

            .copy-btn:hover {
                background-color: var(--raven-teal, #16A085);
                border-color: var(--raven-teal, #16A085);
            }

            .copy-btn.copied {
                background-color: var(--raven-green, #27AE60);
                border-color: var(--raven-green, #27AE60);
            }
        `;
        document.head.appendChild(style);
    }

    /**
     * Handle diagram placeholder loading
     */
    function initDiagramPlaceholders() {
        const placeholders = document.querySelectorAll('.diagram-placeholder');
        placeholders.forEach(placeholder => {
            const diagramName = placeholder.dataset.diagram;
            const img = placeholder.querySelector('img');

            if (img && !img.complete) {
                // Show loading state
                placeholder.style.minHeight = '400px';
                placeholder.style.display = 'flex';
                placeholder.style.alignItems = 'center';
                placeholder.style.justifyContent = 'center';

                // Handle load error
                img.addEventListener('error', () => {
                    const fallback = document.createElement('div');
                    fallback.className = 'diagram-fallback';
                    fallback.innerHTML = `
                        <p style="color: var(--raven-gray); text-align: center;">
                            Diagram: ${diagramName}<br>
                            <em>SVG will be rendered from D2 source</em>
                        </p>
                    `;
                    placeholder.innerHTML = '';
                    placeholder.appendChild(fallback);
                });
            }
        });
    }

    /**
     * Initialize all features
     */
    function init() {
        initSmoothScroll();
        initActiveNav();
        initLazyLoad();
        addCopyButtonStyles();
        initCodeCopy();
        initDiagramPlaceholders();

        console.log('🐦‍⬛ Pedantic Raven - Site loaded');
    }

    // Run when DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();

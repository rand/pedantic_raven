// Theme Toggle
(function() {
    const THEME_KEY = 'pedantic-raven-theme';
    
    // Get saved theme or system preference
    function getInitialTheme() {
        const saved = localStorage.getItem(THEME_KEY);
        if (saved) return saved;
        
        return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    }
    
    // Apply theme
    function applyTheme(theme) {
        document.body.classList.remove('light-theme', 'dark-theme');
        if (theme !== 'system') {
            document.body.classList.add(theme + '-theme');
        }
        localStorage.setItem(THEME_KEY, theme);
    }
    
    // Toggle between light and dark
    function toggleTheme() {
        const currentTheme = localStorage.getItem(THEME_KEY) || 'system';
        const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        
        let newTheme;
        if (currentTheme === 'system') {
            newTheme = systemPrefersDark ? 'light' : 'dark';
        } else if (currentTheme === 'light') {
            newTheme = 'dark';
        } else {
            newTheme = 'light';
        }
        
        applyTheme(newTheme);
    }
    
    // Initialize
    applyTheme(getInitialTheme());
    
    // Add click handler
    const toggleButton = document.querySelector('.theme-toggle');
    if (toggleButton) {
        toggleButton.addEventListener('click', toggleTheme);
    }
    
    // Listen for system theme changes
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
        if (localStorage.getItem(THEME_KEY) === 'system') {
            applyTheme('system');
        }
    });
})();

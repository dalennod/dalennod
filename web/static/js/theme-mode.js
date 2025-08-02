const MODE_KEY = "themeMode";
const getInitialMode = () => {
    const stored = localStorage.getItem(MODE_KEY);
    if (stored === "light" || stored === "dark") {
        return stored;
    }

    const colorScheme = window.matchMedia("(prefers-color-scheme:dark)");
    return colorScheme.matches ? "dark" : "light";
};

const applyMode = (mode) => {
    const rootElement = document.documentElement;
    if (mode === "dark") {
        rootElement.setAttribute("theme-mode", "dark");
    } else {
        rootElement.removeAttribute("theme-mode");
    }
};

const modeToggle = () => {
    const currentMode = document.documentElement.getAttribute("theme-mode") === "dark" ? "dark" : "light";
    const nextMode = currentMode === "dark" ? "light" : "dark";
    applyMode(nextMode);
    localStorage.setItem(MODE_KEY, nextMode);
};

const initThemeMode = (() => {
    const initialMode = getInitialMode();
    applyMode(initialMode);
})();

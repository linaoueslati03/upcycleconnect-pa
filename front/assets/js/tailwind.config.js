// Couleurs et polices de la charte graphique, à charger juste après le CDN Tailwind
tailwind.config = {
  theme: {
    extend: {
      colors: {
        brand: {
          green: "#2E7D32",
          navy: "#1F3864",
          brown: "#8D6E63",
          bg: "#e0dbcc"
        }
      },
      fontFamily: {
        display: ["Poppins", "sans-serif"],
        sans: ["Inter", "sans-serif"]
      }
    }
  }
};

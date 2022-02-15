module.exports = {
  purge: {
    content: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
    enabled: process.env.NODE_ENV == "production",
  },
  //darkMode: 'media', 
  theme: {
    extend: {},
  },
  variants: {
    extend: {},
  },
  plugins: [require("daisyui")],
  daisyui: {
    //themes: ["light", "dark"],
    themes: ["light"],
  },
};

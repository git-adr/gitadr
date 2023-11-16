/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.go.html"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}


/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./templates/**/*.templ"
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('flowbite/plugin'),
    require('tailwindcss'),
    require('autoprefixer'),
  ]
}


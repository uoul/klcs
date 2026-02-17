/** @type {import('tailwindcss').Config} */
import daisyui from "daisyui"

module.exports = {
  content: [
    "./src/**/*.{html,ts}",
  ],
  plugins: [
    daisyui,
  ],
  daisyui: {
    themes: [
      {
        klcs_light: {
           'primary' : '#303030',
           'primary-focus' : '#202020',
           'primary-content' : '#ffffff',

           'secondary' : '#E67817',
           'secondary-focus' : '#ba6012',
           'secondary-content' : '#ffffff',

           'accent' : '#adcd37',
           'accent-focus' : '#88a12b',
           'accent-content' : '#ffffff',

           'neutral' : '#3b424e',
           'neutral-focus' : '#2a2e37',
           'neutral-content' : '#ffffff',

           'base-100' : '#fcfcfc',
           'base-200' : '#dbdbdb',
           'base-300' : '#bbbbbb',
           'base-content' : '#161616',

           'info' : '#1c92f2',
           'success' : '#009485',
           'warning' : '#ff9900',
           'error' : '#ff5724',       
        },
      },
      "dark",
      "cupcake",
    ],
  }
}


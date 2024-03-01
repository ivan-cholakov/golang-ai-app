// tailwind.config.js

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./view/**/*.templ}", "./**/*.templ"],
  dafelist: [
    "border-info",
    "border-b-2",
    "border-white",
    "border-2",
    "bg-base-content",
    "border-base-content",
    "text-black",
    "px-18"
  ],
  theme:{
    extend:{
      fontFamily:{
        sans:["Inter"]
      }
    }
  },
  plugins:[require("daisyui")],
  daisyui: {
    themes:["dark"]
  }
}
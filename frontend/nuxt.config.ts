import NuxtConfiguration from "@nuxt/config";

const config: NuxtConfiguration = {
  mode: "spa",

  /*
   ** Headers of the page
   */
  head: {
    title: "Adventar",
    meta: [{ charset: "utf-8" }, { name: "viewport", content: "width=device-width, initial-scale=1" }],
    link: [{ rel: "icon", type: "image/x-icon", href: "/favicon.ico" }]
  },

  /*
   ** Customize the progress-bar color
   */
  loading: { color: "#fff" },

  /*
   ** Global CSS
   */
  css: [],

  /*
   ** Plugins to load before mounting the App
   */
  plugins: [],

  /*
   ** Nuxt.js modules
   */
  modules: ["@nuxtjs/pwa"],

  /*
   ** Build configuration
   */
  build: {
    babel: {
      ignore: [/\/lib\/grpc\//]
    },
    /*
     ** You can extend webpack config here
     */
    extend(config, ctx) {
      if (ctx.isClient) config.devtool = "#source-map";

      // if (config.module) {
      //   config.module.rules[2] = {};
      // }

      // Run ESLint on save
      // if (ctx.isDev && ctx.isClient && config.module) {
      //   config.module.rules.push({
      //     enforce: "pre",
      //     test: /\.(js|vue)$/,
      //     loader: "eslint-loader",
      //     exclude: /(node_modules)/
      //   });
      // }
    }
  }
};

export default config;

// https://nuxt.com/docs/api/configuration/nuxt-config
import type { NuxtPage } from "nuxt/schema";

export default defineNuxtConfig({
  css: [
    "~/assets/css/tailwind.css",
    "~/assets/css/global.css",
  ],
  plugins: ["~/plugins/echarts.ts"],
  echarts: {
    charts: ["BarChart", "LineChart", "PieChart"],
    components: ["DatasetComponent", "GridComponent", "TooltipComponent"],
    features: ["LabelLayout", "UniversalTransition"],
    renderer: ["svg", "canvas"],
  },
  runtimeConfig: {
    public: {
      API_HOST: process.env.NUXT_PUBLIC_API_HOST, // Memastikan API_HOST tersedia
      WA_HOST: process.env.NUXT_PUBLIC_WA_HOST, // Memastikan API_HOST tersedia
    },
  },
  hooks: {
    "pages:extend"(pages) {
      function setMiddleware(pages: NuxtPage[]) {
        for (const page of pages) {
          if (/* some condition */ true) {
            page.meta ||= {};
            // Note that this will override any middleware set in `definePageMeta` in the page
            page.meta.middleware = ["auth"];
          }
          if (page.children) {
            setMiddleware(page.children);
          }
        }
      }
      setMiddleware(pages);
    },
  },
  app: {
    head: {
      meta: [
        { name: "viewport", content: "width=device-width, initial-scale=1.0" },
      ],
    },
  },
  typescript: {
    typeCheck: true,
  },
  compatibilityDate: "2024-11-01",
  devtools: { enabled: true },
  modules: [
    "@pinia/nuxt",
    "@nuxt/ui",
    "@nuxtjs/leaflet",
    "@i2d/nuxt-pdf-frame",
    "nuxt-echarts",
  ],
  // tailwindcss: {
  //   cssPath: ["~/assets/css/tailwind.css", { injectPosition: "first" }],
  //   config: {},
  //   exposeConfig: true,
  //   viewer: true,
  //   // and more...
  // },
  pinia: {
    storesDirs: ["./stores/**"],
  },
  components: [
    {
      path: "~/components",
      pathPrefix: false,
    },
  ],
});

module.exports = {
  theme: "cosmos",
  title: "Cosmos IBC Account",
  locales: {
    "/": {
      lang: "en-US"
    },
    kr: {
      lang: "kr"
    },
    cn: {
      lang: "cn"
    },
    ru: {
      lang: "ru"
    }
  },
  base: process.env.VUEPRESS_BASE || "/",
  head: [
    ['link', { rel: "apple-touch-icon", sizes: "180x180", href: "/apple-touch-icon.png" }],
    ['link', { rel: "icon", type: "image/png", sizes: "32x32", href: "/favicon-32x32.png" }],
    ['link', { rel: "icon", type: "image/png", sizes: "16x16", href: "/favicon-16x16.png" }],
    ['link', { rel: "manifest", href: "/site.webmanifest" }],
    ['meta', { name: "msapplication-TileColor", content: "#2e3148" }],
    ['meta', { name: "theme-color", content: "#ffffff" }],
    ['link', { rel: "icon", type: "image/svg+xml", href: "/favicon-svg.svg" }],
    ['link', { rel: "apple-touch-icon-precomposed", href: "/apple-touch-icon-precomposed.png" }],
  ],
  themeConfig: {
    repo: "chainapsis/cosmos-sdk-interchain-account",
    docsRepo: "chainapsis/cosmos-sdk-interchain-account",
    docsDir: "docs",
    label: "sdk",
    topbar: {
      banner: false
    },
    sidebar: { 
      auto: false,
      nav: [
        {
          title: "Using the SDK",
          children: [
            {
              title: "Modules",
              directory: true,
              path: "/modules"
            }
          ]
        },
        {
          title: "Test with Starport",
          children: [
            {
              title: "Demo app guide",
              directory: true,
              path: "/starport"
            }
          ]
        }
      ]
    }
  }
};

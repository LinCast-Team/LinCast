const path = require('path');

module.exports = {
  "stories": [
    "../src/**/*.stories.mdx",
    "../src/**/*.stories.@(js|jsx|ts|tsx)"
  ],
  "addons": [
    "@storybook/addon-links",
    "@storybook/addon-essentials",
    "@storybook/addon-postcss",
    "@storybook/addon-a11y",
  ],
  "framework": "@storybook/vue3",
  webpackFinal: async (config, { configType }) => {
    // `configType` has a value of 'DEVELOPMENT' or 'PRODUCTION'
    // You can change the configuration based on that.
    // 'PRODUCTION' is used when building the static version of storybook.

    // Make whatever fine-grained changes you need
    config.module.rules.push(
      {
        test: /\.scss$/,
        use: [
          'style-loader',
          'css-loader',
          {
            loader: 'sass-loader',
            options: { /* Styles automatically applied on stories */
              prependData: `
              @import "@/assets/css/_palette.scss";
              @import "@/assets/css/styles.scss";
              `
            }
          }
        ],
        include: path.resolve(__dirname, '../'),
      },
      // {
      //   test: /.svg$/,
      //   use: ['vue-svg-loader']
      // },
    );

    config.resolve.alias = {
      ...config.resolve.alias,
      '@': path.resolve(__dirname, '../src'),
    }

    // Return the altered config
    return config;
  },
}

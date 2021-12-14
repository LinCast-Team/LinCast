const path = require('path');

module.exports = {
  pluginOptions: {
    'style-resources-loader': {
      preProcessor: 'scss',
      patterns: [ /* Styles that will be included on all the components */
        path.resolve(__dirname, './src/assets/css/_palette.scss'),
        path.resolve(__dirname, './src/assets/css/styles.scss'),
      ],
    },
  },
};

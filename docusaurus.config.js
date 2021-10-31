// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Migrations Documentation',
  tagline: 'Dinosaurs are cool',
  url: 'https://jamillosantos.github.io',
  baseUrl: '/migrations/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',
  organizationName: 'jamillosantos', // Usually your GitHub org/user name.
  projectName: 'migrations', // Usually your repo name.

  presets: [
    [
      '@docusaurus/preset-classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        sidebarPath: require.resolve('./sidebars.js'),
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          routeBasePath: '/',
          // Please change this to your repo.
          editUrl: 'https://github.com/jamillosantos/migrations/edit/main/website/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        title: 'Migrations',
        items: [
          {
            type: 'doc',
            docId: 'introduction',
            position: 'left',
            label: 'Documentation',
          },
          {
            href: 'https://github.com/jamillosantos/migrations',
            label: 'GitHub',
            position: 'right',
          },
          {
            href: 'https://pkg.go.dev/github.com/jamillosantos/migrations',
            label: 'GoDocs',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Documentation',
                to: '/',
              },
              {
                label: 'GitHub',
                href: 'https://github.com/jamillosantos/migrations',
              },
              {
                label: 'GoDocs',
                href: 'https://pkg.go.dev/github.com/jamillosantos/migrations',
              },
            ],
          },
        ],
        copyright: `jamillosantos/migrations is licensed under the MIT license, Inc. Built with Docusaurus.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
};

module.exports = config;

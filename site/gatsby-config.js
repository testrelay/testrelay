module.exports = {
    siteMetadata: {
        siteUrl: "https://testrelay.io",
        title: "Test Relay",
    },
    plugins: [
        "gatsby-plugin-image",
        {
            resolve: "gatsby-plugin-google-analytics",
            options: {
                trackingId: "G-BK7ZYP9EY5",
            },
        },
        "gatsby-plugin-react-helmet",
        "gatsby-plugin-sitemap",
        "gatsby-plugin-sharp",
        "gatsby-transformer-sharp",
        "gatsby-plugin-postcss",
        {
            resolve: "gatsby-source-filesystem",
            options: {
                name: "images",
                path: "./src/images/",
            },
            __key: "images",
        },
        {
            resolve: `gatsby-plugin-manifest`,
            options: {
                name: "TestRelay",
                short_name: "TestRelay",
                start_url: "/",
                background_color: "#1f2937",
                theme_color: "#1f2937",
                display: "standalone",
                icon: "src/images/icon.png",
                crossOrigin: `use-credentials`,
            },
        }
    ],
};

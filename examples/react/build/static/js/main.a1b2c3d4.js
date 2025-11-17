// Main application JavaScript
// This simulates a React bundle with a content hash

console.log('Wayframe React App - Main Bundle');
console.log('Environment:', window.__REACT_ENV__);

// Simulate some React-like functionality
(function() {
    'use strict';

    window.WayframeReact = {
        version: '1.0.0',
        env: window.__REACT_ENV__ || {},

        getEnv: function(key, defaultValue) {
            return this.env[key] || defaultValue;
        },

        apiUrl: function() {
            return this.getEnv('REACT_APP_API_URL', 'http://localhost:8080');
        },

        init: function() {
            console.log('Wayframe React initialized');
            console.log('API URL:', this.apiUrl());
        }
    };

    // Auto-initialize
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', function() {
            window.WayframeReact.init();
        });
    } else {
        window.WayframeReact.init();
    }
})();


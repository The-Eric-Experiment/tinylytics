(function () {
  'use strict';

  class AppWindow extends HTMLElement {
    constructor() {
      super();
      this.attachShadow({ mode: 'open' });
    }

    connectedCallback() {
      const title = this.getAttribute('title') || '';

      // Clone all <link rel="stylesheet"> and <style> tags from document
      const externalStyles = Array.from(
        document.querySelectorAll('link[rel="stylesheet"], style')
      )
        .map((node) => node.cloneNode(true).outerHTML)
        .join('\n');

      this.shadowRoot.innerHTML = `
        ${externalStyles}
        <style>
          :host {
            display: flex;
            flex-direction: column;
          }
          .window {
            display: flex;
            flex-direction: column;
            height: 100%;
            flex: 1;
            min-height: 0;
          }
          .window-padding {
            flex: 1;
            display: flex;
            flex-direction: column;
            min-height: 0;
          }
        </style>
        <div class="window">
          <div class="title-bar">
            <div class="title-bar-text">${title}</div>
          </div>
          <div class="window-padding">
            <slot></slot>
          </div>
        </div>
      `;

    }
  }

  customElements.define('app-window', AppWindow);
})();

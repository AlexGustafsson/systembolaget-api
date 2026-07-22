/* Systembolaget stock card editor card */

const systembolagetStockCardEditor = document.createElement("template");
systembolagetStockCardEditor.innerHTML = `
    <div>
      <ha-card>
        <div class="card-content">
          <ha-input
            id="apiUrl"
            label="API URL"
          >
          </ha-input>
          <ha-input
            id="storeId"
            label="Store ID"
          >
          </ha-input>
          <ha-input
            id="productId"
            label="Product ID"
          >
          </ha-input>
        </div>
      </ha-card>
    </div>
`;

class SystembolagetStockCardEditor extends HTMLElement {
  #apiUrl = "";
  #storeId = "";
  #productId = "";

  constructor() {
    super();
    this.config = {};
    if (!this.shadowRoot) {
      this.attachShadow({ mode: "open" });
      this.shadowRoot.append(
        systembolagetStockCardEditor.content.cloneNode(true),
      );

      this.shadowRoot
        .querySelector("#apiUrl")
        ?.addEventListener("change", (e) => {
          if (e.target.value === this.#apiUrl) {
            return;
          }

          this.#apiUrl = e.target.value;
          this.#onChange();
        });

      this.shadowRoot
        .querySelector("#storeId")
        ?.addEventListener("change", (e) => {
          if (e.target.value === this.#storeId) {
            return;
          }

          this.#storeId = e.target.value;
          this.#onChange();
        });

      this.shadowRoot
        .querySelector("#productId")
        ?.addEventListener("change", (e) => {
          if (e.target.value === this.#productId) {
            return;
          }

          this.#productId = e.target.value;
          this.#onChange();
        });
    }

    this.#update();
  }

  #update() {
    const apiUrl = this.shadowRoot.querySelector("#apiUrl");
    if (apiUrl) {
      apiUrl.value = this.#apiUrl;
    }

    const storeId = this.shadowRoot.querySelector("#storeId");
    if (storeId) {
      storeId.value = this.#storeId;
    }

    const productId = this.shadowRoot.querySelector("#productId");
    if (productId) {
      productId.value = this.#productId;
    }
  }

  #onChange() {
    const event = new Event("config-changed", {
      bubbles: true,
      cancelable: false,
      composed: true,
    });
    event.detail = {
      config: {
        type: "custom:systembolaget-stock-card",
        apiUrl: this.#apiUrl,
        storeId: this.#storeId,
        productId: this.#productId,
      },
    };
    this.dispatchEvent(event);
  }

  setConfig(config) {
    this.#apiUrl = config.apiUrl || "";
    this.#storeId = config.storeId || "";
    this.#productId = config.productId || "";
    this.#update();
  }
}

customElements.define(
  "systembolaget-stock-card-editor",
  SystembolagetStockCardEditor,
);

/* Systembolaget stock card */

const systembolagetStockCard = document.createElement("template");
systembolagetStockCard.innerHTML = `
    <style>
      ha-card {
        height: 100%;
      }

      #image-container {
        position: relative;
        width: 53px;
        height: 104px;
        margin-right: 12px;
      }

      #image-container > img {
        position: absolute;
        top: 0;
        left: 0;
      }

      img {
        object-fit: contain;
      }

      p {
        margin: 0;
      }

      #category {
        font-size: var(--ha-font-size-s);
        font-weight: var(--ha-font-weight-normal);
        line-height: var(--ha-line-height-normal);
        text-transform: uppercase;
        color: var(--secondary-text-color);
      }

      #title {
        font-size: var(--ha-font-size-m);
        font-weight: var(--ha-font-weight-medium);
        line-height: var(--ha-line-height-normal);
        letter-spacing: .1px;
        color: var(--primary-text-color);
      }

      #subtitle {
        font-size: var(--ha-font-size-m);
        font-weight: var(--ha-font-weight-medium);
        line-height: var(--ha-line-height-normal);
        letter-spacing: .1px;
        color: var(--secondary-text-color);
      }

      #number {
        font-size: var(--ha-font-size-s);
        font-weight: var(--ha-font-weight-light);
        line-height: var(--ha-line-height-normal);
        letter-spacing: .1px;
        color: var(--secondary-text-color);
      }

      #details > p {
        font-size: var(--ha-font-size-m);
        font-weight: var(--ha-font-weight-normal);
        line-height: var(--ha-line-height-normal);
        color: var(--secondary-text-color);
      }

      #price {
        flex-grow: 1;
        text-align: right;
        font-weight: var(--ha-font-weight-medium)!important;
      }

      hr {
        border: 0;
        height: 1px;
        background-color: #e8e8e8;
        margin: 6px 0;
      }

      #stock.success::before {
        color: var(--ha-color-fill-success-loud-resting);
      }

      #stock.warn {
        color: var(--ha-color-fill-warn-loud-resting);
      }

      #stock.danger {
        color: var(--ha-color-fill-danger-loud-resting);
      }

      #stock::before {
        content: "";
        display: inline-block;
        width: 0.5em;
        height: 0.5em;
        background-color: currentColor;
        border-radius: 50%;
        margin-right: 0.6em;
        vertical-align: middle;
      }

      .flex { display: flex; }
      .flex-column { flex-direction: column; }
      .flex-grow { flex-grow: 1; }
      .items-center { align-items: center; }
      .justify-between { justify-content: space-between; }
      .space-12 > :not(:last-child) { margin-right: 12px }
      .padding-10-20-0-20 { padding: 10px 20px 0 20px; }
      .padding-0-20-10-20 { padding: 0 20px 10px 20px; }
    </style>
    <ha-card>
      <div class="flex padding-10-20-0-20">
        <div id="image-container">
          <img id="thumbnail" alt="Product thumbnail" width="53" height="104" />
          <img id="image" alt="Product image" width="53" height="104" loading="lazy" decoding="async" />
        </div>
        <div class="flex flex-column flex-grow">
          <p id="category"></p>
          <p id="title"></p>
          <p id="subtitle"></p>
          <p id="number"></p>
          <div class="flex items-center justify-between space-12">
            <p id="country"></p>
            <p id="volume"></p>
            <p id="alcohol"></p>
            <p id="price"></p>
          </div>
        </div>
      </div>
      <hr />
      <div class="flex justify-between padding-0-20-10-20">
        <div>
          <p>Antal i butik</p>
          <p id="stock"></p>
        </div>
        <div class="flex space-12">
          <div>
            <p>Sektion</p>
            <p id="section"></p>
          </div>
          <div>
            <p>Hylla</p>
            <p id="shelf"></p>
          </div>
        <div>
      </div>
    </ha-card>
`;

class SystembolagetStockCard extends HTMLElement {
  /**
   * @type {{
   *  stock: number,
   *  shelf: string,
   *  category: string | undefined,
   *  title: string | undefined,
   *  subtitle: string | undefined,
   *  number: string | undefined,
   *  country: string | undefined,
   *  volume: string | undefined,
   *  alcoholPercentage: number | undefined,
   *  price: number | undefined,
   *  thumbnail: string | undefined,
   *  imageUrl: string | undefined,
   * } | undefined}
   */
  #product = undefined;

  #apiUrl = "";
  #productId = "";
  #storeId = "";

  constructor() {
    super();
    if (!this.shadowRoot) {
      this.attachShadow({ mode: "open" });
      this.shadowRoot.append(systembolagetStockCard.content.cloneNode(true));
    }

    this.#update();
  }

  // Whenever the state changes, a new `hass` object is set. Use this to
  // update your content.
  set hass(_hass) {
    this.#update();
  }

  #update() {
    if (!this.shadowRoot) {
      return;
    }

    if (!this.#product && this.#apiUrl && this.#storeId && this.#productId) {
      fetch(
        `${this.#apiUrl}/stores/${this.#storeId}/products/${this.#productId}`,
      )
        .then((res) => {
          if (res.status !== 200) {
            throw new Error(
              `Failed to fetch product: unexpected status code ${res.status}`,
            );
          }

          return res.json();
        })
        .then((product) => {
          this.#product = product;
          this.#setValues();
        })
        .catch((error) => {
          console.error(error);
        });
    }
  }

  #setValues() {
    if (!this.shadowRoot || !this.#product) {
      return;
    }

    const thumbnail = this.shadowRoot.querySelector("#thumbnail");
    if (thumbnail) {
      thumbnail.setAttribute(
        "src",
        `data:image/png;base64,${this.#product.thumbnail}`,
      );
    }

    const image = this.shadowRoot.querySelector("#image");
    if (image && this.#product.imageUrl) {
      image.setAttribute("src", this.#product.imageUrl);
      image.addEventListener(
        "error",
        () => {
          image?.parentNode?.removeChild(image);
        },
        { once: true },
      );
      if (thumbnail) {
        image.addEventListener(
          "load",
          () => {
            thumbnail?.parentNode?.removeChild(thumbnail);
          },
          { once: true },
        );
      }
    }

    const category = this.shadowRoot.querySelector("#category");
    if (category) {
      category.textContent = this.#product.category || "";
    }

    const title = this.shadowRoot.querySelector("#title");
    if (title) {
      title.textContent = this.#product.title || "";
    }

    const subtitle = this.shadowRoot.querySelector("#subtitle");
    if (subtitle) {
      subtitle.textContent = this.#product.subtitle || "";
    }

    const number = this.shadowRoot.querySelector("#number");
    if (number) {
      number.textContent = this.#product.number
        ? `Nr ${this.#product.number}`
        : "";
    }

    const country = this.shadowRoot.querySelector("#country");
    if (country) {
      country.textContent = this.#product.country || "";
    }

    const volume = this.shadowRoot.querySelector("#volume");
    if (volume) {
      volume.textContent = this.#product.volume || "";
    }

    const alcoholPercentage = this.shadowRoot.querySelector("#alcohol");
    if (alcoholPercentage) {
      alcoholPercentage.textContent = this.#product.alcoholPercentage
        ? `${this.#product.alcoholPercentage} % vol.`
        : "Alkoholfri";
    }

    const price = this.shadowRoot.querySelector("#price");
    if (price) {
      price.textContent = this.#product.price
        ? this.#product.price.toFixed(2).replace(".", ":")
        : "";
    }

    const stock = this.shadowRoot.querySelector("#stock");
    if (stock) {
      if (this.#product.stock === 0) {
        stock.className = "danger";
      } else if (this.#product.stock < 5) {
        stock.className = "warn";
      } else {
        stock.className = "success";
      }
      stock.textContent = `${this.#product.stock} st` || "";
    }

    const segments = this.#product.shelf.split("-").filter((x) => x);
    const section = this.shadowRoot.querySelector("#section");
    if (section && segments.length >= 1) {
      section.textContent = segments[0];
    }
    const shelf = this.shadowRoot.querySelector("#shelf");
    if (shelf && segments.length >= 2) {
      shelf.textContent = segments[1];
    }
  }

  // The user supplied configuration. Throw an exception and Home Assistant
  // will render an error card.
  setConfig(config) {
    if (!config.apiUrl || !config.storeId || !config.productId) {
      throw new Error("Invalid settings");
    }

    this.#apiUrl = config.apiUrl;
    this.#storeId = config.storeId;
    this.#productId = config.productId;
  }

  // The height of your card. Home Assistant uses this to automatically
  // distribute all cards over the available columns in masonry view
  getCardSize() {
    return 3;
  }

  // The rules for sizing your card in the grid in sections view
  getGridOptions() {
    return {
      columns: 12,
      min_columns: 9,
      min_rows: 3,
      max_rows: 3,
    };
  }

  static getConfigElement() {
    return document.createElement("systembolaget-stock-card-editor");
  }

  static getStubConfig() {
    return {};
  }
}

customElements.define("systembolaget-stock-card", SystembolagetStockCard);

/** Expose in Lovelace */

window.customCards.push({
  type: "systembolaget-stock-card",
  name: "Systembolaget stock card",
  preview: false,
  description: "Shows the current stock of a product in Systembolaget",
});

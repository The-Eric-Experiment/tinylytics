// World Map Web Component (Shadow DOM)
// Self-contained component using D3.js for rendering
// Usage:
//   <world-map geo-url="/static/worldmap.json" height="400px">
//     <country iso="US" count="300"></country>
//     <country iso="AU" count="150"></country>
//     <country iso="BR" count="75"></country>
//   </world-map>

import { colord } from 'colord';

(function () {
    'use strict';

    // ISO_COUNTRIES mapping - 3-letter (GeoJSON) to 2-letter (ISO) codes
    const ISO_3_TO_2 = {
        "AFG": "AF", "ALA": "AX", "ALB": "AL", "DZA": "DZ", "ASM": "AS", "AND": "AD",
        "AGO": "AO", "AIA": "AI", "ATA": "AQ", "ATG": "AG", "ARG": "AR", "ARM": "AM",
        "ABW": "AW", "AUS": "AU", "AUT": "AT", "AZE": "AZ", "BHS": "BS", "BHR": "BH",
        "BGD": "BD", "BRB": "BB", "BLR": "BY", "BEL": "BE", "BLZ": "BZ", "BEN": "BJ",
        "BMU": "BM", "BTN": "BT", "BOL": "BO", "BIH": "BA", "BWA": "BW", "BVT": "BV",
        "BRA": "BR", "VGB": "VG", "IOT": "IO", "BRN": "BN", "BGR": "BG", "BFA": "BF",
        "BDI": "BI", "KHM": "KH", "CMR": "CM", "CAN": "CA", "CPV": "CV", "CYM": "KY",
        "CAF": "CF", "TCD": "TD", "CHL": "CL", "CHN": "CN", "HKG": "HK", "MAC": "MO",
        "CXR": "CX", "CCK": "CC", "COL": "CO", "COM": "KM", "COG": "CG", "COD": "CD",
        "COK": "CK", "CRI": "CR", "CIV": "CI", "HRV": "HR", "CUB": "CU", "CYP": "CY",
        "CZE": "CZ", "DNK": "DK", "DJI": "DJ", "DMA": "DM", "DOM": "DO", "ECU": "EC",
        "EGY": "EG", "SLV": "SV", "GNQ": "GQ", "ERI": "ER", "EST": "EE", "ETH": "ET",
        "FLK": "FK", "FRO": "FO", "FJI": "FJ", "FIN": "FI", "FRA": "FR", "GUF": "GF",
        "PYF": "PF", "ATF": "TF", "GAB": "GA", "GMB": "GM", "GEO": "GE", "DEU": "DE",
        "GHA": "GH", "GIB": "GI", "GRC": "GR", "GRL": "GL", "GRD": "GD", "GLP": "GP",
        "GUM": "GU", "GTM": "GT", "GGY": "GG", "GIN": "GN", "GNB": "GW", "GUY": "GY",
        "HTI": "HT", "HMD": "HM", "VAT": "VA", "HND": "HN", "HUN": "HU", "ISL": "IS",
        "IND": "IN", "IDN": "ID", "IRN": "IR", "IRQ": "IQ", "IRL": "IE", "IMN": "IM",
        "ISR": "IL", "ITA": "IT", "JAM": "JM", "JPN": "JP", "JEY": "JE", "JOR": "JO",
        "KAZ": "KZ", "KEN": "KE", "KIR": "KI", "PRK": "KP", "KOR": "KR", "KWT": "KW",
        "KGZ": "KG", "LAO": "LA", "LVA": "LV", "LBN": "LB", "LSO": "LS", "LBR": "LR",
        "LBY": "LY", "LIE": "LI", "LTU": "LT", "LUX": "LU", "MKD": "MK", "MDG": "MG",
        "MWI": "MW", "MYS": "MY", "MDV": "MV", "MLI": "ML", "MLT": "MT", "MHL": "MH",
        "MTQ": "MQ", "MRT": "MR", "MUS": "MU", "MYT": "YT", "MEX": "MX", "FSM": "FM",
        "MDA": "MD", "MCO": "MC", "MNG": "MN", "MNE": "ME", "MSR": "MS", "MAR": "MA",
        "MOZ": "MZ", "MMR": "MM", "NAM": "NA", "NRU": "NR", "NPL": "NP", "NLD": "NL",
        "NCL": "NC", "NZL": "NZ", "NIC": "NI", "NER": "NE", "NGA": "NG", "NIU": "NU",
        "NFK": "NF", "MNP": "MP", "NOR": "NO", "OMN": "OM", "PAK": "PK", "PLW": "PW",
        "PSE": "PS", "PAN": "PA", "PNG": "PG", "PRY": "PY", "PER": "PE", "PHL": "PH",
        "PCN": "PN", "POL": "PL", "PRT": "PT", "PRI": "PR", "QAT": "QA", "REU": "RE",
        "ROU": "RO", "RUS": "RU", "RWA": "RW", "BLM": "BL", "SHN": "SH", "KNA": "KN",
        "LCA": "LC", "MAF": "MF", "SPM": "PM", "VCT": "VC", "WSM": "WS", "SMR": "SM",
        "STP": "ST", "SAU": "SA", "SEN": "SN", "SRB": "RS", "SYC": "SC", "SLE": "SL",
        "SGP": "SG", "SVK": "SK", "SVN": "SI", "SLB": "SB", "SOM": "SO", "ZAF": "ZA",
        "SGS": "GS", "SSD": "SS", "ESP": "ES", "LKA": "LK", "SDN": "SD", "SUR": "SR",
        "SJM": "SJ", "SWZ": "SZ", "SWE": "SE", "CHE": "CH", "SYR": "SY", "TWN": "TW",
        "TJK": "TJ", "TZA": "TZ", "THA": "TH", "TLS": "TL", "TGO": "TG", "TKL": "TK",
        "TON": "TO", "TTO": "TT", "TUN": "TN", "TUR": "TR", "TKM": "TM", "TCA": "TC",
        "TUV": "TV", "UGA": "UG", "UKR": "UA", "ARE": "AE", "GBR": "GB", "USA": "US",
        "UMI": "UM", "URY": "UY", "UZB": "UZ", "VUT": "VU", "VEN": "VE", "VNM": "VN",
        "VIR": "VI", "WLF": "WF", "ESH": "EH", "YEM": "YE", "ZMB": "ZM", "ZWE": "ZW",
        // Special territories
        "XKX": "XK", "XNC": "XN", "XSO": "XS"
    };


    // Default colors matching the React version
    const DEFAULT_COLORS = {
        baseColor: '#01FF00',
        fillColor: '#008000',
        strokeColor: '#008000',
        hoverColor: '#01FF00'
    };

    // Note: This component requires ac-colors via CDN:
    // <script src="https://cdn.jsdelivr.net/npm/ac-colors@1/ac-colors.min.js"></script>

    class WorldMapComponent extends HTMLElement {
        constructor() {
            super();
            this._shadow = this.attachShadow({ mode: 'open' });
            this._data = new Map();
            this._colors = { ...DEFAULT_COLORS };
            this._geoUrl = '/worldmap.json';
            this._geoData = null;
            this._rendered = false;

            this._shadow.innerHTML = `
                <style>
                    :host {
                        display: block;
                        position: relative;
                        width: 100%;
                        height: 100%;
                        min-height: 0;
                    }
                    .container {
                        width: 100%;
                        height: 100%;
                        display: flex;
                        align-items: center;
                        justify-content: center;
                        overflow: hidden;
                    }
                    svg {
                        display: block;
                        max-width: 100%;
                        max-height: 100%;
                        width: auto;
                        height: auto;
                    }
                    .tooltip {
                        position: fixed;
                        background: #ffffff;
                        border: 2px solid #000000;
                        padding: 5px 10px;
                        pointer-events: none;
                        opacity: 0;
                        font-family: system-ui, -apple-system, sans-serif;
                        font-size: 12px;
                        z-index: 10000;
                        box-shadow: 2px 2px 4px rgba(0,0,0,0.3);
                        transition: opacity 0.15s;
                        white-space: nowrap;
                    }
                    .error {
                        color: red;
                        padding: 20px;
                        text-align: center;
                    }
                    path {
                        cursor: pointer;
                        outline: none;
                    }
                </style>
                <div class="container">
                    <svg viewBox="0 0 800 400" preserveAspectRatio="xMidYMid meet"></svg>
                </div>
                <div class="tooltip"></div>
            `;

            this._container = this._shadow.querySelector('.container');
            this._svg = this._shadow.querySelector('svg');
            this._tooltip = this._shadow.querySelector('.tooltip');
        }

        static get observedAttributes() {
            return ['geo-url', 'base-color', 'fill-color', 'stroke-color', 'hover-color', 'height'];
        }

        attributeChangedCallback(name, oldValue, newValue) {
            if (oldValue === newValue) return;
            switch (name) {
                case 'geo-url':
                    this._geoUrl = newValue;
                    this._geoData = null;
                    if (this._rendered) this._render();
                    break;
                case 'base-color':
                    this._colors.baseColor = newValue;
                    if (this._rendered) this._updateColors();
                    break;
                case 'fill-color':
                    this._colors.fillColor = newValue;
                    if (this._rendered) this._updateColors();
                    break;
                case 'stroke-color':
                    this._colors.strokeColor = newValue;
                    if (this._rendered) this._updateColors();
                    break;
                case 'hover-color':
                    this._colors.hoverColor = newValue;
                    break;
                case 'height':
                    this.style.height = newValue;
                    break;
            }
        }

        connectedCallback() {
            this._parseCountryElements();

            this._observer = new MutationObserver((mutations) => {
                const hasCountryChanges = mutations.some(m => {
                    for (const n of m.addedNodes) if (n.nodeName === 'COUNTRY') return true;
                    for (const n of m.removedNodes) if (n.nodeName === 'COUNTRY') return true;
                    if (m.type === 'attributes' && m.target.nodeName === 'COUNTRY') return true;
                    return false;
                });
                if (hasCountryChanges) {
                    this._parseCountryElements();
                    if (this._rendered) this._updateColors();
                }
            });
            this._observer.observe(this, {
                childList: true,
                subtree: true,
                attributes: true,
                attributeFilter: ['iso', 'count']
            });

            // Wait for D3 (required)
            this._waitForD3()
                .then(() => this._render())
                .catch(err => {
                    console.error('World Map: D3.js not available', err);
                    this._container.innerHTML = '<div class="error">D3.js not available. Please include: &lt;script src="https://d3js.org/d3.v7.min.js"&gt;&lt;/script&gt;</div>';
                });
        }

        disconnectedCallback() {
            if (this._observer) this._observer.disconnect();
        }

        _parseCountryElements() {
            this._data.clear();
            this.querySelectorAll('country').forEach(el => {
                let iso = (el.getAttribute('iso') || '').toUpperCase();
                const count = parseInt(el.getAttribute('count') || '0', 10);
                if (iso.length === 3 && ISO_3_TO_2[iso]) iso = ISO_3_TO_2[iso];
                if (iso && !isNaN(count)) this._data.set(iso, count);
            });
        }

        _waitForD3() {
            return new Promise((resolve, reject) => {
                if (typeof d3 !== 'undefined') return resolve();
                let attempts = 0;
                const check = setInterval(() => {
                    attempts++;
                    if (typeof d3 !== 'undefined') { clearInterval(check); resolve(); }
                    else if (attempts >= 50) { clearInterval(check); reject(new Error('D3 timeout')); }
                }, 100);
            });
        }


        _getFillColor(geoJsonId) {
            const iso2 = ISO_3_TO_2[geoJsonId];
            const count = this._data.get(iso2);
            if (count === undefined) return this._colors.fillColor;

            const lightenAmount = 0.4 * (1.0 - count / 100);

            return colord(this._colors.baseColor)
                .lighten(lightenAmount)
                .toHex();
        }

        _getOpacity(geoJsonId) {
            return 1;
        }

        _getTooltipText(geoJsonId, feature) {
            const iso2 = ISO_3_TO_2[geoJsonId];
            // Use name from GeoJSON properties if available, fallback to ISO code
            const name = (feature?.properties?.name) || iso2 || geoJsonId;
            return `${name}: ${this._data.get(iso2) || 0} sessions`;
        }

        _showTooltip(text, event) {
            this._tooltip.textContent = text;
            this._tooltip.style.opacity = '1';
            this._tooltip.style.left = (event.clientX + 12) + 'px';
            this._tooltip.style.top = (event.clientY - 12) + 'px';
        }

        _hideTooltip() {
            this._tooltip.style.opacity = '0';
        }

        _render() {
            this._rendered = true;
            d3.select(this._svg).selectAll('*').remove();

            const width = 800;
            const height = 400;

            // Update viewBox
            this._svg.setAttribute('viewBox', `0 0 ${width} ${height}`);

            // Use geoNaturalEarth1 projection - curved, rounded appearance
            const projection = d3.geoNaturalEarth1();
            const pathGen = d3.geoPath().projection(projection);

            const renderPaths = (geo) => {
                if (!geo?.features) {
                    this._container.innerHTML = '<div class="error">Invalid GeoJSON</div>';
                    return;
                }

                // Filter out Antarctica to prevent it from affecting bounds and taking up space
                const filteredFeatures = geo.features.filter(feature => {
                    const iso2 = ISO_3_TO_2[feature.id];
                    return iso2 !== 'AQ';
                });

                // Create a new FeatureCollection without Antarctica for bounds calculation
                const filteredGeo = {
                    type: 'FeatureCollection',
                    features: filteredFeatures
                };

                // Fit the projection to the filtered GeoJSON bounds (without Antarctica)
                projection.fitSize([width, height], filteredGeo);

                d3.select(this._svg).selectAll('path')
                    .data(filteredFeatures)
                    .enter()
                    .append('path')
                    .attr('d', pathGen)
                    .attr('fill', d => this._getFillColor(d.id) || 'transparent')
                    .attr('stroke', this._colors.strokeColor)
                    .attr('stroke-width', 0.5)
                    .attr('opacity', d => this._getOpacity(d.id))
                    .on('mouseover', (e, d) => {
                        if (ISO_3_TO_2[d.id] === 'AQ') return;
                        this._showTooltip(this._getTooltipText(d.id, d), e);
                        d3.select(e.currentTarget).transition().duration(100)
                            .attr('fill', this._colors.hoverColor);
                    })
                    .on('mousemove', (e) => {
                        this._tooltip.style.left = (e.clientX + 12) + 'px';
                        this._tooltip.style.top = (e.clientY - 12) + 'px';
                    })
                    .on('mouseout', (e, d) => {
                        this._hideTooltip();
                        d3.select(e.currentTarget).transition().duration(100)
                            .attr('fill', this._getFillColor(d.id) || 'transparent');
                    });
            };

            if (this._geoData) {
                renderPaths(this._geoData);
            } else {
                d3.json(this._geoUrl)
                    .then(geo => { this._geoData = geo; renderPaths(geo); })
                    .catch(err => {
                        console.error('World Map: GeoJSON load error', err);
                        this._container.innerHTML = `<div class="error">Error loading ${this._geoUrl}</div>`;
                    });
            }
        }

        _updateColors() {
            d3.select(this._svg).selectAll('path')
                .transition().duration(200)
                .attr('fill', d => this._getFillColor(d.id) || 'transparent')
                .attr('stroke', this._colors.strokeColor);
        }

        // Public API
        setData(data) {
            this._data.clear();
            data.forEach(item => {
                let iso = (item.iso || item.value || '').toUpperCase();
                if (iso.length === 3 && ISO_3_TO_2[iso]) iso = ISO_3_TO_2[iso];
                if (iso) this._data.set(iso, item.count || 0);
            });
            if (this._rendered) this._updateColors();
        }

        getData() { return new Map(this._data); }

        setCountryCount(iso, count) {
            iso = iso.toUpperCase();
            if (iso.length === 3 && ISO_3_TO_2[iso]) iso = ISO_3_TO_2[iso];
            this._data.set(iso, count);
            if (this._rendered) this._updateColors();
        }
    }

    if (!customElements.get('world-map')) {
        customElements.define('world-map', WorldMapComponent);
    }
})();
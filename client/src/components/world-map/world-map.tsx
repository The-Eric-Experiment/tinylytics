import { colord } from "colord";
import React, { FunctionComponent, useState } from "react";
import { ComposableMap, Geographies, Geography } from "react-simple-maps";
import ReactTooltip from "react-tooltip";
import styled from "styled-components";
import { AnalyticsData } from "../../api/types";
import { COUNTRIES, ISO_COUNTRIES } from "../../constants/countries";

interface WorldMapProps {
  data?: AnalyticsData[];
}

export const WorldMap: FunctionComponent<WorldMapProps> = ({ data }) => {
  const [tooltip, setTooltip] = useState<string | null>(null);
  const colors = {
    baseColor: "#01FF00",
    fillColor: "#008000",
    strokeColor: "#008000",
    hoverColor: "#01FF00",
  };

  function getFillColor(code: string) {
    const iso = ISO_COUNTRIES[code];
    if (iso === "AQ") return;
    const country = data?.find(({ value }) => value === iso);

    if (!country) {
      return colors.fillColor;
    }

    return colord(colors.baseColor)
      .lighten(0.4 * (1.0 - country.count / 100))
      .toHex();
  }

  function getOpacity(code: string) {
    const iso = ISO_COUNTRIES[code];
    return iso === "AQ" ? 0 : 1;
  }

  function handleHover(code: string) {
    const iso = ISO_COUNTRIES[code];
    if (iso === "AQ") return;
    const country = data?.find(({ value }) => value === iso);
    setTooltip(`${COUNTRIES[iso]}: ${country?.count || 0} sessions`);
  }

  return (
    <Container data-tip="" data-for="world-map-tooltip">
      <ComposableMap
        projection="geoMercator"
        projectionConfig={{
          center: [0, 40],
          scale: 120,
        }}
      >
        <Geographies geography="/worldmap.json">
          {({ geographies }) => {
            return geographies.map((geo) => {
              return (
                <Geography
                  key={geo.rsmKey}
                  geography={geo}
                  fill={getFillColor(geo.id)}
                  stroke={colors.strokeColor}
                  opacity={getOpacity(geo.id)}
                  style={{
                    default: { outline: "none" },
                    hover: { outline: "none", fill: colors.hoverColor },
                    pressed: { outline: "none" },
                  }}
                  onMouseOver={() => handleHover(geo.id)}
                  onMouseOut={() => setTooltip(null)}
                />
              );
            });
          }}
        </Geographies>
      </ComposableMap>
      <ReactTooltip id="world-map-tooltip">{tooltip}</ReactTooltip>
    </Container>
  );
};

const Container = styled.div`
  position: relative;
  overflow: hidden;
`;

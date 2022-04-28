import React, { FunctionComponent } from "react";
import Flag from "react-world-flags";
import { Panel } from "react95";
import styled from "styled-components";
import { useCountries } from "../../api/analytics";
import { AnalyticsData } from "../../api/types";
import { COUNTRIES } from "../../constants/countries";
import AnalyticsTable from "../shared/analytics-table";
import { LANDSCAPE } from "../shared/page-layout";
import { SuspenseWindow } from "../shared/suspense-window";
import {
  TableWidgetContainer,
  TableWidgetWrapperProps,
} from "../shared/table-widget";
import { WorldMap } from "../world-map/world-map";

interface CountriesProps extends Omit<TableWidgetWrapperProps, "title"> {
  domain: string;
}

const CountriesContent: FunctionComponent<CountriesProps> = ({
  domain,
  ...props
}) => {
  const { data } = useCountries(domain, props.filters);

  const hasPreviousFilters =
    data?.previousFilters && data.previousFilters.length > 0;

  const renderIcon = (item: AnalyticsData) => (
    <Flag code={item.value} width={16} />
  );

  const formatName = (name: string) => {
    return COUNTRIES[name] || name;
  };

  return (
    <Container>
      {hasPreviousFilters && <div>{data.previousFilters.join(", ")}</div>}
      <Row>
        <MapPanel variant="well">
          <WorldMap data={data?.items} />
        </MapPanel>
        <TableWidgetContainer>
          <AnalyticsTable
            {...props}
            data={data}
            renderIcon={renderIcon}
            formatName={formatName}
            filterPrimary="c"
            showSelfWhenEmpty
          />
        </TableWidgetContainer>
      </Row>
    </Container>
  );
};

export const Countries: FunctionComponent<CountriesProps> = (props) => {
  return (
    <SuspenseWindow title="Countries">
      <CountriesContent {...props} />
    </SuspenseWindow>
  );
};

const Container = styled.div`
  display: flex;
  flex-direction: column;
`;

const Row = styled.div`
  gap: 20px;
  display: flex;
  flex-direction: column;

  @media all and (min-width: ${LANDSCAPE}) {
    flex-direction: row;
  }

  & > div {
    flex: 1;
  }
`;

const MapPanel = styled(Panel)`
  background-color: #000080;
`;

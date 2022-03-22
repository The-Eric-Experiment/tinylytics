import React from "react";
import { useSummaries } from "../../api/analytics";
import { Filters } from "../../api/types";
import { Card, GridItem } from "../shared/page-layout";

export interface SummaryProps {
  domain: string;
  filters: Filters;
}

export function Summary({ domain, filters }: SummaryProps) {
  const { data, error } = useSummaries(domain, filters);

  const renderCardContent = (content: () => React.ReactNode) => {
    if (error) {
      return <div>Whoops...</div>;
    }

    if (!data) {
      return <div>Loading...</div>;
    }

    return <div>{content()}</div>;
  };

  return (
    <>
      <GridItem take={1}>
        <Card>
          <b>Sessions: </b> {renderCardContent(() => data?.sessions)}
        </Card>
      </GridItem>
      <GridItem take={1}>
        <Card>
          <b>PageViews: </b> {renderCardContent(() => data?.pageViews)}
        </Card>
      </GridItem>
      <GridItem take={1}>
        <Card>
          <b>Bounce rate: </b>
        </Card>
      </GridItem>
      <GridItem take={1}>
        <Card>
          <b>Avg Session Duration: </b>
        </Card>
      </GridItem>
    </>
  );
}

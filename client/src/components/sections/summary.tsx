import React from "react";
import { useSummaries } from "../../api/analytics";
import { Filters } from "../../api/types";
import { Card, GridItem } from "../shared/page-layout";

export interface SummaryProps {
  domain: string;
  filters: Filters;
}

export function Summary({ domain, filters }: SummaryProps) {
  const { data } = useSummaries(domain, filters);

  return (
    <>
      <GridItem take={1}>
        <Card>
          <b>Sessions: </b> {data!.sessions}{" "}
        </Card>
      </GridItem>
      <GridItem take={1}>
        <Card>
          <b>PageViews: </b> {data!.pageViews}{" "}
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

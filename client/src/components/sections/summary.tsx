import React from "react";
import { useSummaries } from "../../api/analytics";
import { Filters } from "../../api/types";

export interface SummaryProps {
  domain: string;
  filters: Filters;
}

export function Summary({ domain, filters }: SummaryProps) {
  const { data } = useSummaries(domain, filters);

  return (
    <ul>
      <li>
        <b>Sessions: </b> {data!.sessions}{" "}
      </li>
      <li>
        <b>PageViews: </b> {data!.pageViews}{" "}
      </li>
    </ul>
  );
}

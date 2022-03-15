import React from "react";
import { useSummaries } from "../../api/analytics";

export interface SummaryProps {
  domain: string;
}

export function Summary({ domain }: SummaryProps) {
  const { error, data } = useSummaries(domain);

  if (error)
    return <div> "An error has occurred: " + {(error as any).message} </div>;

  if (!data) return <div>Loading...</div>;

  return (
    <ul>
      <li>
        <b>Sessions: </b> {data?.sessions}{" "}
      </li>
      <li>
        <b>PageViews: </b> {data?.pageViews}{" "}
      </li>
    </ul>
  );
}

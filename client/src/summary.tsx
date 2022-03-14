import React from "react";
import { useQuery } from "react-query";
import { getSummaries } from "./api/analytics";

export function Widget() {
  const { isLoading, error, data, isFetching } = useQuery("summaries", () =>
    getSummaries("oldavista.com")
  );

  if (isLoading) return <div>Loading...</div>;

  if (error)
    return <div> "An error has occurred: " + {(error as any).message} </div>;

  return (
    <ul>
      <li>
        <b>Sessions: </b> {data?.sessions}{" "}
      </li>
      <li>
        <b>PageViews: </b> {data?.pageViews}{" "}
      </li>
      <li>{isFetching ? "Updating..." : ""}</li>
    </ul>
  );
}

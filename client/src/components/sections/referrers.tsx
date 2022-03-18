import React, { FunctionComponent } from "react";
import { useReferrers } from "../../api/analytics";
import { AnalyticsData, Filters } from "../../api/types";
import { useAnalyticsFilters } from "../../hooks/analytics-data-hooks";

export interface ReferrersProps {
  domain: string;
  filters: Filters;
  onFilter(filters: Partial<Filters>): void;
}

export const Referrers: FunctionComponent<ReferrersProps> = ({
  domain,
  filters,
  onFilter,
}) => {
  const { data } = useReferrers(domain, filters);
  const { getLabel, showAsLink, updateFilter, selectedFilter } =
    useAnalyticsFilters(filters, { name: "r", version: "rv" }, onFilter);

  const getLocalLabel = (data: AnalyticsData) => {
    const label = getLabel(data);
    if (!label) {
      return "(other)";
    }
    return label;
  };

  return (
    <div>
      {selectedFilter && <div>{selectedFilter}</div>}
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Count</th>
          </tr>
        </thead>
        <tbody>
          {data!.items.map((item) => (
            <tr key={item.name + getLabel(item)}>
              <td>
                {showAsLink(item) ? (
                  <a onClick={updateFilter(item)} type="button" href="">
                    {getLocalLabel(item)}
                  </a>
                ) : (
                  getLocalLabel(item)
                )}
              </td>
              <td>{item.count}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

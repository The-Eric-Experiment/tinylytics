import React, { FunctionComponent } from "react";
import { useBrowsers } from "../../api/analytics";
import { Filters } from "../../api/types";
import { useAnalyticsFilters } from "../../hooks/analytics-data-hooks";

interface BrowsersProps {
  domain: string;
  filters: Filters;
  onFilter(filters: Partial<Filters>): void;
}

export const Browsers: FunctionComponent<BrowsersProps> = ({
  domain,
  filters,
  onFilter,
}) => {
  const { data } = useBrowsers(domain, filters);
  const { getLabel, showAsLink, updateFilter, selectedFilter } =
    useAnalyticsFilters(filters, { name: "b", version: "bv" }, onFilter);

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
                    {getLabel(item)}
                  </a>
                ) : (
                  getLabel(item)
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

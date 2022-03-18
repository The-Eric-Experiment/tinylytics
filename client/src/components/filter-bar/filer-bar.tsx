import React, { FunctionComponent } from "react";
import { Filters, Periods } from "../../api/types";
import { FILTER_NAMES, PERIOD_NAMES } from "../../constants/filters";
import { Pill } from "../shared/page-layout";

type FilterBarProps = {
  filters: Filters;
  removeFilter(key: keyof Filters): void;
  onFilter(filters: Partial<Filters>): void;
};

export const FilterBar: FunctionComponent<FilterBarProps> = ({
  filters,
  removeFilter,
  onFilter,
}) => {
  const removeFilterClick = (key: keyof Filters) => () => removeFilter(key);
  const changePeriod = (value: string) => {
    onFilter({ p: value as Periods });
  };
  return (
    <>
      <select onChange={(item) => changePeriod(item.target.value)}>
        {(Object.keys(PERIOD_NAMES) as Array<Periods>).map((key) => {
          return (
            <option key={key} value={key} selected={key === filters.p}>
              {PERIOD_NAMES[key]}
            </option>
          );
        })}
      </select>
      {(Object.keys(filters) as Array<keyof Filters>)
        .filter((key) => key !== "p")
        .map((key) => {
          return (
            <Pill>
              {FILTER_NAMES[key]} is {filters[key]}
              <button key={key} onClick={removeFilterClick(key)}>
                x
              </button>
            </Pill>
          );
        })}
    </>
  );
};

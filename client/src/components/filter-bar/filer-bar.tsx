import React, { FunctionComponent, useMemo } from "react";
import { Bar, Panel, Select, SelectOption, Toolbar } from "react95";
import { Filters, Periods } from "../../api/types";
import {
  DEPENDANT_FILTERS,
  FILTER_NAMES,
  PERIODS,
  SHOW_AS_SAME_FILTER,
  SHOW_PREVIOUS_FILTER_IF_EMPTY,
} from "../../constants/filters";
import { Pill } from "../shared/page-layout";
import { WebsiteSelector } from "./website-selector";

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
  const changePeriod = (
    _: React.FormEvent,
    nextSelection: SelectOption<Periods>
  ) => {
    onFilter({ p: nextSelection.value });
  };

  const presentFilters = useMemo(
    () =>
      (Object.keys(filters) as Array<keyof Filters>)
        .filter((key) => key !== "p")
        .reduce((acc, key) => {
          const showAsSame = SHOW_AS_SAME_FILTER.find((sas) =>
            sas.includes(key)
          );
          if (!showAsSame) {
            return [...acc, key];
          }

          const inAccIndex = acc.findIndex((a) => showAsSame.includes(a));
          if (inAccIndex === -1) {
            return [...acc, key];
          }

          const indexCurrent = showAsSame.indexOf(key);
          const indexExisting = showAsSame.indexOf(acc[inAccIndex]);

          if (indexCurrent > indexExisting) {
            return acc.map((k, i) => {
              if (i === inAccIndex) {
                return key;
              }
              return k;
            });
          }

          return [...acc, key];
        }, [] as Array<keyof Filters>),
    [filters]
  );

  const getName = (key: keyof Filters) => {
    const item = filters[key];
    if (!item && SHOW_PREVIOUS_FILTER_IF_EMPTY.includes(key)) {
      const previous = (
        Object.keys(DEPENDANT_FILTERS) as Array<keyof Filters>
      ).find((k) => DEPENDANT_FILTERS[k]?.includes(key));
      return filters[previous || key];
    }

    if (!item) {
      return "(unknown)";
    }

    return item;
  };

  return (
    <Panel style={{ width: "100%" }}>
      <Toolbar>
        <Bar size={35} />
        <Select
          defaultValue={filters.p}
          options={PERIODS}
          menuMaxHeight={160}
          width={160}
          onChange={
            changePeriod as React.ChangeEventHandler<HTMLInputElement> &
              ((
                e: React.FormEvent<Element>,
                nextSelection: SelectOption<Periods>
              ) => void)
          }
        />

        <WebsiteSelector />

        <Bar size={35} />
        {presentFilters.map((key) => {
          return (
            <Pill key={key}>
              {FILTER_NAMES[key]} is {getName(key)}
              <button key={key} onClick={removeFilterClick(key)}>
                x
              </button>
            </Pill>
          );
        })}
      </Toolbar>
    </Panel>
  );
};

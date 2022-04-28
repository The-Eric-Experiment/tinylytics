import { Duration } from "luxon";
import React from "react";
import { Cutout } from "react95";
import styled from "styled-components";
import { useSummaries } from "../../api/analytics";
import { Filters } from "../../api/types";
import { BaseWindow } from "../shared/base-window";
import { GridItemX1 } from "../shared/page-layout";

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

  const getAvgSessionDuration = () => {
    return Duration.fromObject({ seconds: data?.avgSessionDuration }).toFormat(
      "hh:mm:ss"
    );
  };

  return (
    <>
      <GridItemX1>
        <BaseWindow title="Sessions">
          <CutoutStyled>{renderCardContent(() => data?.sessions)}</CutoutStyled>
        </BaseWindow>
      </GridItemX1>
      <GridItemX1>
        <BaseWindow title="Page Views">
          <CutoutStyled>
            {renderCardContent(() => data?.pageViews)}
          </CutoutStyled>
        </BaseWindow>
      </GridItemX1>
      <GridItemX1>
        <BaseWindow title="Bounce Rate">
          <CutoutStyled>
            {renderCardContent(() => `${data?.bounceRate}%`)}
          </CutoutStyled>
        </BaseWindow>
      </GridItemX1>
      <GridItemX1>
        <BaseWindow title="Avg. Session Duration">
          <CutoutStyled>
            {renderCardContent(getAvgSessionDuration)}
          </CutoutStyled>
        </BaseWindow>
      </GridItemX1>
    </>
  );
}

const CutoutStyled = styled(Cutout)`
  background: #ffffff;
  padding: 0 16px;
  text-align: right;
  font-size: 42px;
`;

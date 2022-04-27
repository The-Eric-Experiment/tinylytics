import styled from "styled-components";

export const LANDSCAPE = "667px";
export const TABLET = "760px";

export const PageLayout = styled.div`
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 16px;
  padding: 16px;
  width: 100%;
  max-width: 1600px;
`;

export const PageHeader = styled.header`
  display: flex;
  flex-direction: row;
  gap: 16px;
`;

export const PageGrid = styled.section`
  display: grid;
  grid-template-rows: auto;
  grid-template-columns: repeat(2, 1fr);
  grid-gap: 16px;

  @media all and (min-width: ${LANDSCAPE}) {
    grid-template-columns: repeat(4, 1fr);
  }
`;

export const GridItemX1 = styled.div`
  grid-column: span 1;
  @media all and (min-width: ${LANDSCAPE}) {
    grid-column: span 1;
  }
  display: flex;
  flex-direction: column;
`;

export const GridItemX2 = styled.div`
  grid-column: span 2;
  @media all and (min-width: ${LANDSCAPE}) {
    grid-column: span 2;
  }
  display: flex;
  flex-direction: column;
`;

export const GridItemX4 = styled.div`
  grid-column: span 2;
  @media all and (min-width: ${LANDSCAPE}) {
    grid-column: span 4;
  }
  display: flex;
  flex-direction: column;
`;

export const Card = styled.div`
  display: flex;
  flex-direction: column;
  border: 1px solid black;
  padding: 16px;
  height: 100%;
`;

export const Pill = styled.div`
  border-radius: 5px;
  border: 1px solid black;
`;

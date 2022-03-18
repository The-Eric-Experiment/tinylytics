import styled from "styled-components";

export const PageLayout = styled.div`
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 16px;
  padding: 16px;
`;

export const PageHeader = styled.header`
  display: flex;
  flex-direction: row;
  gap: 16px;
`;

export const PageGrid = styled.section`
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  grid-template-rows: auto;
  grid-gap: 16px;
`;

export const GridItem = styled.div<{ take: number }>`
  grid-column: span ${(props) => props.take};
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

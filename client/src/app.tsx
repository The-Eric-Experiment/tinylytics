import React from "react";
import { Route, Routes } from "react-router-dom";
import styled from "styled-components";
import { AnalyticsPage } from "./pages/analytics-page";

function App() {
  return (
    <Container>
      <Routes>
        <Route path="/" element={<AnalyticsPage />} />
      </Routes>
    </Container>
  );
}

const Container = styled.div`
  display: flex;
  flex-direction: row;
  align-items: flex-start;
  justify-content: space-between;
`;

export default App;

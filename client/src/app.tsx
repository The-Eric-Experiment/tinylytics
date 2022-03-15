import React from "react";
import styled from "styled-components";
import { Browsers } from "./components/sections/browsers";
import { Countries } from "./components/sections/countries";
import { OSs } from "./components/sections/os";
import { Summary } from "./components/sections/summary";

function App() {
  const domain = "oldavista.com";
  return (
    <Container>
      <Summary domain={domain} />
      <Browsers domain={domain} />
      <OSs domain={domain} />
      <Countries domain={domain} />
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

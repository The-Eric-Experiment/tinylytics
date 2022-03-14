import React from "react";
import { QueryClient, QueryClientProvider } from "react-query";
import { Widget } from "./summary";

const queryClient = new QueryClient();

function App() {
  return (
    <div>
      <QueryClientProvider client={queryClient}>
        <Widget />
      </QueryClientProvider>
    </div>
  );
}

export default App;

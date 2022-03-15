import React, { FunctionComponent } from "react";
import { useBrowsers } from "../../api/analytics";

interface BrowsersProps {
  domain: string;
}

export const Browsers: FunctionComponent<BrowsersProps> = ({ domain }) => {
  const { error, data } = useBrowsers(domain);

  if (error) return <div> "An error has occurred: " + {error.message} </div>;

  if (!data) return <div>Loading...</div>;

  return (
    <table>
      <thead>
        <tr>
          <th>Name</th>
          <th>Count</th>
        </tr>
      </thead>
      <tbody>
        {data.items.map((item) => (
          <tr key={item.name}>
            <td>{item.name}</td>
            <td>{item.count}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

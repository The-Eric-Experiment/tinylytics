import React, { FunctionComponent } from "react";
import { useOSs } from "../../api/analytics";

interface OSsProps {
  domain: string;
}

export const OSs: FunctionComponent<OSsProps> = ({ domain }) => {
  const { error, data } = useOSs(domain);

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

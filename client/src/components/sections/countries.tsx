import React, { FunctionComponent } from "react";
import { useCountries } from "../../api/analytics";

export interface CountriesProps {
  domain: string;
}

export const Countries: FunctionComponent<CountriesProps> = ({ domain }) => {
  const { error, data } = useCountries(domain);

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

import React, { FunctionComponent, useMemo } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useWebsites } from "../../api/analytics";
import { Selector } from "../shared/selector";

type WebsiteSelectorProps = {};

export const WebsiteSelector: FunctionComponent<
  WebsiteSelectorProps
> = ({}) => {
  const websites = useWebsites()?.data || [];
  const navigate = useNavigate();
  const params = useParams();

  const selected = websites.find((o) => o.domain === params.domain)!;
  console.log(selected);

  const options = useMemo(() => {
    return websites.map((o) => ({
      label: o.title,
      value: o.domain,
    }));
  }, [websites]);

  const changeSite = (value: string) => {
    navigate(`/${value}`);
  };

  return (
    <Selector
      options={options}
      selectedValue={selected?.domain}
      onChange={changeSite}
    />
  );
};

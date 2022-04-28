import React, { FunctionComponent, useMemo } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Select, SelectOption } from "react95";
import { useWebsites } from "../../api/analytics";

type WebsiteSelectorProps = {};

export const WebsiteSelector: FunctionComponent<WebsiteSelectorProps> = () => {
  const analyticsWebsites = useWebsites();
  const websites = useMemo(
    () => analyticsWebsites?.data || [],
    [analyticsWebsites]
  );
  const navigate = useNavigate();
  const params = useParams();

  const selected = websites.find((o) => o.domain === params.domain)!;

  const options: SelectOption<string>[] = useMemo(() => {
    return websites.map((o) => ({
      label: o.title,
      value: o.domain,
    }));
  }, [websites]);

  const changeSite = (
    _: React.FormEvent,
    nextSelection: SelectOption<string>
  ) => {
    navigate(`/${nextSelection.value}`);
  };

  return (
    <Select
      defaultValue={selected.domain}
      options={options}
      menuMaxHeight={160}
      width={160}
      onChange={
        changeSite as React.ChangeEventHandler<HTMLInputElement> &
          ((
            e: React.FormEvent<Element>,
            nextSelection: SelectOption<string>
          ) => void)
      }
    />
  );
};

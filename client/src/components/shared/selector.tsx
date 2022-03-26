import React, { useCallback, useEffect, useMemo, useState } from "react";
import styled from "styled-components";

export type Option<TValue> = {
  value: TValue;
  label: string;
};

export type Separator = {
  isSeparator: true;
};

export type Item<TValue> = Separator | Option<TValue>;

function isSeparator<TValue>(input: Item<TValue>): input is Separator {
  return !!(input as Separator).isSeparator;
}

const DropDownContainer = styled.div`
  width: 10.5em;
`;

const DropDownHeader = styled.div`
  margin-bottom: 0.8em;
  padding: 0.4em 2em 0.4em 1em;
  border: 1px solid #000;
  font-weight: 500;
  cursor: pointer;
`;

const DropDownListContainer = styled.div`
  position: absolute;
  z-index: 100;
  min-width: 200px;
  box-shadow: 0 2px 3px rgba(0, 0, 0, 0.15);
`;

const DropDownList = styled.ul`
  padding: 0;
  margin: 0;
  background: #ffffff;
  border: 1px solid #000;
  box-sizing: border-box;
  font-weight: 500;
`;

const Separator = styled.li`
  list-style: none;
  border-bottom: 1px solid #000;
`;

const ListItem = styled.li`
  list-style: none;
  padding: 12px 20px;
  cursor: pointer;

  &:hover {
    background-color: #ddd;
  }
`;

type SelectorProps<TValue> = {
  options: Array<Item<TValue>>;
  selectedValue: TValue;
  onChange(item: TValue): void;
};

export function Selector<TValue>({
  selectedValue,
  options,
  onChange,
}: SelectorProps<TValue>) {
  const [isOpen, setIsOpen] = useState(false);

  const selected = useMemo(
    () => (options as Option<TValue>[]).find((o) => o.value === selectedValue),
    [selectedValue]
  );

  const close = useCallback(() => {
    setIsOpen(false);
  }, []);

  useEffect(() => {
    if (isOpen) {
      document.addEventListener("click", close);
    } else {
      document.removeEventListener("click", close);
    }
    return () => document.removeEventListener("click", close);
  }, [isOpen]);

  const toggling = () => setIsOpen(!isOpen);

  const onOptionClicked = (item: Option<TValue>) => () => {
    onChange(item.value);
    setIsOpen(false);
  };

  const renderItem = (option: Item<TValue>, index: number) => {
    if (isSeparator(option)) {
      return <Separator key={index} />;
    }

    return (
      <ListItem onClick={onOptionClicked(option)} key={index}>
        {option.label}
      </ListItem>
    );
  };

  return (
    <DropDownContainer>
      <DropDownHeader onClick={toggling}>{selected?.label}</DropDownHeader>
      {isOpen && (
        <DropDownListContainer>
          <DropDownList>{options.map(renderItem)}</DropDownList>
        </DropDownListContainer>
      )}
    </DropDownContainer>
  );
}

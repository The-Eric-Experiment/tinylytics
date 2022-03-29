import React, { useCallback, useEffect, useMemo, useState } from "react";
import styled, { css } from "styled-components";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowDown,
  faChevronDown,
  faCoffee,
} from "@fortawesome/free-solid-svg-icons";

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
      <DropDownHeader onClick={toggling} isOpen={isOpen}>
        {selected?.label} <FontAwesomeIcon icon={faChevronDown} />
      </DropDownHeader>
      {isOpen && (
        <DropDownListContainer>
          <DropDownList>{options.map(renderItem)}</DropDownList>
        </DropDownListContainer>
      )}
    </DropDownContainer>
  );
}

const DropDownContainer = styled.div``;

const DropDownHeader = styled.div<{ isOpen: boolean }>`
  margin-bottom: 14px;
  padding: 8px 32px 8px 16px;
  border: 1px solid #000;
  font-weight: 500;
  cursor: pointer;
  position: relative;

  > svg {
    position: absolute;
    right: 8px;
    transform: rotate(0);
    transition: transform ease-out 0.2s;
    ${(props) =>
      props.isOpen &&
      css`
        transform: rotate(180deg);
      `};
  }
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

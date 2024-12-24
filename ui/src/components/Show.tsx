import React from "react";

export type ShowProps = {
  when: boolean;
  children: React.ReactNode;
  fallback?: React.ReactNode;
};

const Show = ({ when, children, fallback }: ShowProps) => {
  return when ? children : fallback;
};

export default Show;

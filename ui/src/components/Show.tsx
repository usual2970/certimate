import React from "react";

const Show = ({
  when,
  children,
  fallback,
}: {
  when: boolean;
  children: React.ReactNode;
  fallback?: React.ReactNode;
}) => {
  return when ? children : fallback;
};

export default Show;

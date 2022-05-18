export const defineSpecColor = (spec) => {
  if (!spec.startedAt) {
    return ["", "white"];
  }
  if (!spec.finishedAt) {
    return ["running", "yellow-300"];
  }

  const status = {
    passed: "green-600",
    failed: "red-600",
    unknown: "white",
  };

  return status[spec.status] || status.unknown;
};

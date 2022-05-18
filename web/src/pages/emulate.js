/* eslint-disable no-restricted-globals */

import { CreateSessionForm } from "../components/form/CreateSession";

export const Emulate = () => {
  return (
    <div className="h-full pt-4 sm:pt-12">
      <CreateSessionForm history={history} />
    </div>
  );
};

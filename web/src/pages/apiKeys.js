import { useCallback, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { displayTimestamp } from "../format/displayDate";
import Loading from "../components/atoms/Loading";
import Spinner from "../components/atoms/Spinner";
import { ShowAlert } from "../components/atoms/Alert";
import useFetch from "use-http";
import { endpoints } from "../api";

export const ApiKeys = () => {
  const [apiKeys, setApiKeys] = useState([]);
  const {
    get,
    del: deleteApiKey,
    loading,
    response,
    error,
  } = useFetch(endpoints.apiKeys);

  const [deletionKeyId, setDeletionKey] = useState();

  const fetchKeys = useCallback(async () => {
    const keys = await get();
    if (keys && !error) {
      setApiKeys(keys);
    }
  }, [get, error]);

  useEffect(() => {
    fetchKeys();
  }, [fetchKeys]);

  const currentDate = new Date();
  const currentTimestamp = currentDate.valueOf();

  const onDelete = useCallback(
    async (e) => {
      e.preventDefault();
      const id = e.target.value;
      setDeletionKey(id);
      await deleteApiKey(id);
      setApiKeys(apiKeys.filter((key) => key.id !== id));
    },
    [deleteApiKey, setApiKeys, apiKeys]
  );

  return (
    <div className="h-full pt-4 sm:pt-12">
      <div className="max-w-6xl px-4 mx-auto mt-8">
        {error &&
          ShowAlert(response?.data?.errors?.join("; ") || response?.data)}
        <button className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-8 my-4 rounded">
          <Link to={`/createApiKey`}>Create api key</Link>
        </button>
        {loading ? (
          <Loading />
        ) : apiKeys?.length > 0 ? (
          <table className="table-auto border-collapse border border-blue-400 w-full">
            <thead className="space-x-1">
              <tr className="bg-blue-600 px-auto py-auto">
                <th className="w-1/2 border border-blue-400">
                  <span className="text-gray-100 font-semibold">Name</span>
                </th>
                <th className="w-1/3 border border-blue-400">
                  <span className="text-gray-100 font-semibold">ExpireAt</span>
                </th>
                <th className="w-1/8 border border-blue-400">
                  <span className="text-gray-100 font-semibold"></span>
                </th>
              </tr>
            </thead>
            <tbody className="bg-gray-200">
              {apiKeys
                ?.sort((a, b) => b.expireAt - a.expireAt)
                .map((apiKey) => (
                  <tr key={apiKey.id} className="bg-white">
                    <td className="font-semibold border border-blue-400">
                      {apiKey.name}
                    </td>
                    <td className="border border-blue-400 max-w-0">
                      <p
                        className={
                          currentTimestamp > apiKey.expireAt
                            ? "bg-red-200"
                            : "bg-white"
                        }
                      >
                        {displayTimestamp(apiKey.expireAt)}
                      </p>
                    </td>
                    <td className="border border-blue-400">
                      <button
                        className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-2 rounded w-48"
                        title="Delete api key"
                        onClick={onDelete}
                        value={apiKey.id}
                      >
                        {loading && deletionKeyId === apiKey.id ? (
                          <Spinner />
                        ) : (
                          `Delete`
                        )}
                      </button>
                    </td>
                  </tr>
                ))}
            </tbody>
          </table>
        ) : (
          NoApiKeys()
        )}
      </div>
    </div>
  );
};

const NoApiKeys = () => {
  return (
    <div className="max-w-6xl px-4 mx-auto mt-8">No API Keys available.</div>
  );
};

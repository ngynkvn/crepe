import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { ExternalLink, Loader2 } from "lucide-react";

interface Repository {
  id: string;
  repo_name: string;
  repo_type: string;
  url: string;
  updated_at: string;
  file_count: number;
}

export function RepositoryList() {
  const [repositories, setRepositories] = useState<Repository[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchRepositories = async () => {
      try {
        const response = await fetch("/api/repositories");
        if (!response.ok) {
          throw new Error("Failed to fetch repositories");
        }
        const data = await response.json();
        setRepositories(data);
      } catch (err) {
        setError(`An error occurred while fetching repositories: ${err}`);
      } finally {
        setIsLoading(false);
      }
    };

    fetchRepositories();
  }, []);

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center text-red-500 p-4">
        <p>{error}</p>
      </div>
    );
  }

  return (
    <Card className="w-full max-w-3xl mx-auto">
      <CardHeader>
        <CardTitle>Repositories</CardTitle>
      </CardHeader>
      <CardContent>
        <ScrollArea className="h-[400px] pr-4">
          <ul className="space-y-4">
            {repositories.map((repo) => (
              <li key={repo.id} className="border-b pb-4 last:border-b-0">
                <h3 className="font-semibold text-lg">
                  {repo.repo_name}{" "}
                  <span className="text-muted-foreground text-sm">
                    ({repo.repo_type})
                  </span>
                </h3>
                <div className="flex items-center space-x-2 text-sm text-muted-foreground">
                  <a
                    href={repo.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center hover:text-primary"
                  >
                    <ExternalLink className="h-4 w-4 mr-1" />
                    {repo.url}
                  </a>
                </div>
                <p className="text-sm text-muted-foreground mt-1">
                  Last indexed: {new Date(repo.updated_at).toLocaleString()}
                </p>
                <p className="text-sm text-muted-foreground mt-1">
                  {repo.file_count} files
                </p>
              </li>
            ))}
          </ul>
        </ScrollArea>
      </CardContent>
    </Card>
  );
}

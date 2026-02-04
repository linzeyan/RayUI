import { describe, it, expect, beforeEach } from "vitest";
import { useProfileStore } from "./profileStore";
import { mockAppApi } from "@/test/wails-mock";

describe("profileStore", () => {
  beforeEach(() => {
    useProfileStore.setState({
      profiles: [],
      speedResults: new Map(),
      selectedIds: new Set(),
      filterSubId: "all",
      searchQuery: "",
      loading: false,
    });
  });

  describe("loadProfiles", () => {
    it("loads profiles from backend", async () => {
      const mockProfiles = [
        { id: "1", remarks: "Server 1" },
        { id: "2", remarks: "Server 2" },
      ];
      mockAppApi.GetProfiles.mockResolvedValue(mockProfiles);

      await useProfileStore.getState().loadProfiles();

      expect(mockAppApi.GetProfiles).toHaveBeenCalledWith("");
      expect(useProfileStore.getState().profiles).toEqual(mockProfiles);
    });

    it("converts 'all' filter to empty string for API", async () => {
      mockAppApi.GetProfiles.mockResolvedValue([]);

      await useProfileStore.getState().loadProfiles("all");

      expect(mockAppApi.GetProfiles).toHaveBeenCalledWith("");
    });

    it("passes subscription ID directly to API", async () => {
      mockAppApi.GetProfiles.mockResolvedValue([]);

      await useProfileStore.getState().loadProfiles("sub-123");

      expect(mockAppApi.GetProfiles).toHaveBeenCalledWith("sub-123");
    });

    it("sets loading state correctly", async () => {
      mockAppApi.GetProfiles.mockImplementation(
        () =>
          new Promise((resolve) => {
            expect(useProfileStore.getState().loading).toBe(true);
            resolve([]);
          })
      );

      await useProfileStore.getState().loadProfiles();

      expect(useProfileStore.getState().loading).toBe(false);
    });
  });

  describe("addProfile", () => {
    it("adds profile and reloads", async () => {
      const profile = { id: "1", remarks: "New Server" };
      mockAppApi.AddProfile.mockResolvedValue(undefined);
      mockAppApi.GetProfiles.mockResolvedValue([profile]);

      await useProfileStore.getState().addProfile(profile as never);

      expect(mockAppApi.AddProfile).toHaveBeenCalledWith(profile);
      expect(mockAppApi.GetProfiles).toHaveBeenCalled();
    });
  });

  describe("updateProfile", () => {
    it("updates profile and reloads", async () => {
      const profile = { id: "1", remarks: "Updated Server" };
      mockAppApi.UpdateProfile.mockResolvedValue(undefined);
      mockAppApi.GetProfiles.mockResolvedValue([profile]);

      await useProfileStore.getState().updateProfile(profile as never);

      expect(mockAppApi.UpdateProfile).toHaveBeenCalledWith(profile);
      expect(mockAppApi.GetProfiles).toHaveBeenCalled();
    });
  });

  describe("deleteProfiles", () => {
    it("deletes profiles and removes from selection", async () => {
      useProfileStore.setState({
        selectedIds: new Set(["1", "2", "3"]),
      });
      mockAppApi.DeleteProfiles.mockResolvedValue(undefined);
      mockAppApi.GetProfiles.mockResolvedValue([]);

      await useProfileStore.getState().deleteProfiles(["1", "2"]);

      expect(mockAppApi.DeleteProfiles).toHaveBeenCalledWith(["1", "2"]);
      expect(useProfileStore.getState().selectedIds).toEqual(new Set(["3"]));
    });
  });

  describe("setActiveProfile", () => {
    it("calls backend to set active profile", async () => {
      mockAppApi.SetActiveProfile.mockResolvedValue(undefined);

      await useProfileStore.getState().setActiveProfile("profile-1");

      expect(mockAppApi.SetActiveProfile).toHaveBeenCalledWith("profile-1");
    });
  });

  describe("importFromText", () => {
    it("imports profiles from text and returns count", async () => {
      mockAppApi.ImportFromText.mockResolvedValue(5);
      mockAppApi.GetProfiles.mockResolvedValue([]);

      const count = await useProfileStore.getState().importFromText("vmess://...");

      expect(mockAppApi.ImportFromText).toHaveBeenCalledWith("vmess://...");
      expect(count).toBe(5);
    });
  });

  describe("exportShareLink", () => {
    it("exports profile share link", async () => {
      mockAppApi.ExportShareLink.mockResolvedValue("vmess://exported");

      const link = await useProfileStore.getState().exportShareLink("profile-1");

      expect(mockAppApi.ExportShareLink).toHaveBeenCalledWith("profile-1");
      expect(link).toBe("vmess://exported");
    });
  });

  describe("testProfiles", () => {
    it("tests specific profiles and updates results", async () => {
      const results = [
        { profileId: "1", latency: 100, speed: 1000 },
        { profileId: "2", latency: 200, speed: 2000 },
      ];
      mockAppApi.TestProfiles.mockResolvedValue(results);

      await useProfileStore.getState().testProfiles(["1", "2"]);

      expect(mockAppApi.TestProfiles).toHaveBeenCalledWith(["1", "2"]);
      const speedResults = useProfileStore.getState().speedResults;
      expect(speedResults.get("1")).toEqual(results[0]);
      expect(speedResults.get("2")).toEqual(results[1]);
    });
  });

  describe("testAllProfiles", () => {
    it("tests all profiles and updates results", async () => {
      const results = [{ profileId: "1", latency: 100, speed: 1000 }];
      mockAppApi.TestAllProfiles.mockResolvedValue(results);

      await useProfileStore.getState().testAllProfiles();

      expect(mockAppApi.TestAllProfiles).toHaveBeenCalled();
      expect(useProfileStore.getState().speedResults.get("1")).toEqual(results[0]);
    });
  });

  describe("selection management", () => {
    it("setSelectedIds updates selection", () => {
      useProfileStore.getState().setSelectedIds(new Set(["1", "2"]));
      expect(useProfileStore.getState().selectedIds).toEqual(new Set(["1", "2"]));
    });

    it("toggleSelected adds ID if not present", () => {
      useProfileStore.getState().toggleSelected("1");
      expect(useProfileStore.getState().selectedIds).toEqual(new Set(["1"]));
    });

    it("toggleSelected removes ID if present", () => {
      useProfileStore.setState({ selectedIds: new Set(["1", "2"]) });
      useProfileStore.getState().toggleSelected("1");
      expect(useProfileStore.getState().selectedIds).toEqual(new Set(["2"]));
    });
  });

  describe("filter management", () => {
    it("setFilterSubId updates filter", () => {
      useProfileStore.getState().setFilterSubId("sub-123");
      expect(useProfileStore.getState().filterSubId).toBe("sub-123");
    });

    it("setSearchQuery updates search", () => {
      useProfileStore.getState().setSearchQuery("japan");
      expect(useProfileStore.getState().searchQuery).toBe("japan");
    });
  });
});

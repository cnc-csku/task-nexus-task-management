import MainContainer from "@/components/templates/MainContainer";
import Header from "@/components/ui/Header";
import ProjectsTable from "./_components/ProjectsTable";

export default function ProjectPage() {
  return (
    <MainContainer>
      <Header>Projects</Header>
      <ProjectsTable />
    </MainContainer>
  );
}

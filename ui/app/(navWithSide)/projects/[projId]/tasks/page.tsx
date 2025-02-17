import Header from '@/components/ui/Header'
import React from 'react'
import TasksTable from './_components/TasksTable'

export default function ProjectTaskPage() {
  return (
    <div>
        <Header>
            Tasks
        </Header>
        <TasksTable />
    </div>
  )
}

import { Modal, ModalBody, ModalContent, ModalHeader } from "@nextui-org/modal";

interface SprintGoalModalProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
}

export default function SprintGoalModal({ isOpen, onOpenChange }: SprintGoalModalProps) {
  return (
    <Modal
      isOpen={isOpen}
      onOpenChange={onOpenChange}
    >
      <ModalContent className="py-3">
        {(onClose) => (
          <>
            <ModalHeader>Sprint Goal</ModalHeader>
            <ModalBody className="mx-3">
                <ul className="list-disc">
                    <li>Goal 1</li>
                    <li>Goal 2</li>
                    <li>Goal 3</li> 
                </ul>
            </ModalBody>
          </>
        )}
      </ModalContent>
    </Modal>
  );
}
